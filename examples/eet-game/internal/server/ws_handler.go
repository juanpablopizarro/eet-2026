package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"eet-game/internal/game"
	"eet-game/internal/protocol"
)

const (
	stateInterval = time.Second / 60 // ~60 fps broadcast
	writeTimeout  = 5 * time.Second
	readTimeout   = 30 * time.Second
)

// wsHandler handles one WebSocket connection, creating and running a dedicated
// game instance for the connected client.
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // allow cross-origin in dev
	})
	if err != nil {
		log.Printf("ws accept error: %v", err)
		return
	}
	defer conn.CloseNow()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Event channel so the game can push events to the writer goroutine.
	eventCh := make(chan protocol.GameEventMessage, 64)

	g := game.New(0, func(ev protocol.GameEventMessage) {
		select {
		case eventCh <- ev:
		default:
		}
	})

	// Run the game loop in a goroutine.
	go g.Run(ctx)

	// Writer goroutine — broadcasts state frames and events.
	go func() {
		ticker := time.NewTicker(stateInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case ev := <-eventCh:
				wctx, wcancel := context.WithTimeout(ctx, writeTimeout)
				if err := wsjson.Write(wctx, conn, ev); err != nil {
					wcancel()
					cancel()
					return
				}
				wcancel()
			case <-ticker.C:
				snap := g.Snapshot()
				wctx, wcancel := context.WithTimeout(ctx, writeTimeout)
				if err := wsjson.Write(wctx, conn, snap); err != nil {
					wcancel()
					cancel()
					return
				}
				wcancel()
			}
		}
	}()

	// Reader loop — processes player input messages.
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("ws read error: %v", err)
			return
		}

		var msg struct {
			Type   protocol.MessageType  `json:"type"`
			Action protocol.InputAction  `json:"action"`
		}
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		if msg.Type == protocol.TypeInput {
			g.SendInput(msg.Action)
		}
	}
}
