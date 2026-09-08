package game

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"eet-game/internal/protocol"
	"eet-game/pkg/vec2"
)

// State represents a high-level game state machine state.
type State string

const (
	StateMenu     State = "menu"
	StatePlaying  State = "playing"
	StatePaused   State = "paused"
	StateGameOver State = "game_over"
	StateVictory  State = "victory"
)

const (
	tickRate          = 60           // fps
	tickInterval      = time.Second / tickRate
	alienAnimInterval = 0.6          // seconds between animation frame toggles
	ufoSpawnInterval  = 22.0         // seconds between UFO appearances
	bombInterval      = 1.4          // seconds between alien bomb drops
	maxPlayerBullets  = 1            // player can have at most 1 bullet in flight
)

// EventCallback is called by Game when notable events occur.
type EventCallback func(ev protocol.GameEventMessage)

// Game holds the complete authoritative state of one game session.
type Game struct {
	mu sync.Mutex

	state   State
	player  Player
	aliens  []*Alien
	bullets []*Bullet
	shields []*Shield
	ufo     UFO

	score   int
	hiScore int
	level   int

	alienTimer float64
	alienAnimT float64
	ufoTimer   float64
	bombTimer  float64

	alienDirX float64 // +1 or -1

	inputCh chan protocol.InputAction
	onEvent EventCallback

	rng *rand.Rand
}

// New creates and returns a Game in the menu state.
func New(hiScore int, onEvent EventCallback) *Game {
	g := &Game{
		state:   StateMenu,
		hiScore: hiScore,
		inputCh: make(chan protocol.InputAction, 32),
		onEvent: onEvent,
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return g
}

// SendInput enqueues a player input action (non-blocking).
func (g *Game) SendInput(action protocol.InputAction) {
	select {
	case g.inputCh <- action:
	default:
	}
}

// Run starts the game loop and blocks until ctx is cancelled.
func (g *Game) Run(ctx context.Context) {
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()
	lastTick := time.Now()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			dt := now.Sub(lastTick).Seconds()
			lastTick = now
			g.tick(dt)
		}
	}
}

// Snapshot returns a serialisable snapshot of the current game state.
func (g *Game) Snapshot() protocol.GameStateMessage {
	g.mu.Lock()
	defer g.mu.Unlock()

	msg := protocol.GameStateMessage{
		Type:    protocol.TypeState,
		Score:   g.score,
		HiScore: g.hiScore,
		Level:   g.level,
		State:   string(g.state),
		Player: protocol.PlayerState{
			X:          g.player.Pos.X,
			Y:          g.player.Pos.Y,
			Lives:      g.player.Lives,
			Invincible: g.player.Invincible,
		},
		UFO: protocol.UFOState{
			X:      g.ufo.Pos.X,
			Y:      g.ufo.Pos.Y,
			Active: g.ufo.Active,
		},
	}
	for _, a := range g.aliens {
		msg.Aliens = append(msg.Aliens, protocol.AlienState{
			ID:    a.ID,
			X:     a.Pos.X,
			Y:     a.Pos.Y,
			Type:  int(a.Type),
			Alive: a.Alive,
			Frame: a.Frame,
		})
	}
	for _, b := range g.bullets {
		if b.Alive {
			msg.Bullets = append(msg.Bullets, protocol.BulletState{
				X:     b.Pos.X,
				Y:     b.Pos.Y,
				Owner: string(b.Owner),
			})
		}
	}
	for _, s := range g.shields {
		msg.Shields = append(msg.Shields, protocol.ShieldState{
			X:  s.Pos.X,
			Y:  s.Pos.Y,
			HP: s.HP,
		})
	}
	return msg
}

// ---- Internal tick -------------------------------------------------------

func (g *Game) tick(dt float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.processInputs()
	if g.state == StatePlaying {
		g.update(dt)
	}
}

func (g *Game) processInputs() {
	for {
		select {
		case action := <-g.inputCh:
			g.handleAction(action)
		default:
			return
		}
	}
}

func (g *Game) handleAction(action protocol.InputAction) {
	switch action {
	case protocol.ActionStart:
		if g.state == StateMenu || g.state == StateGameOver || g.state == StateVictory {
			g.startGame()
		}
	case protocol.ActionRestart:
		g.startGame()
	case protocol.ActionPause:
		if g.state == StatePlaying {
			g.state = StatePaused
		} else if g.state == StatePaused {
			g.state = StatePlaying
		}
	case protocol.ActionMoveLeft:
		if g.state == StatePlaying {
			g.player.Pos.X -= PlayerSpeed * (1.0 / tickRate)
			if g.player.Pos.X < 0 {
				g.player.Pos.X = 0
			}
		}
	case protocol.ActionMoveRight:
		if g.state == StatePlaying {
			g.player.Pos.X += PlayerSpeed * (1.0 / tickRate)
			if g.player.Pos.X+PlayerWidth > CanvasWidth {
				g.player.Pos.X = CanvasWidth - PlayerWidth
			}
		}
	case protocol.ActionShoot:
		if g.state == StatePlaying {
			g.tryPlayerShoot()
		}
	}
}

func (g *Game) startGame() {
	g.level = 1
	g.score = 0
	g.startLevel()
}

func (g *Game) startLevel() {
	g.player = NewPlayer()
	g.aliens = buildAlienFormation(g.level)
	g.shields = buildShields()
	g.bullets = nil
	g.ufo = UFO{}
	g.alienDirX = 1
	g.alienTimer = 0
	g.alienAnimT = 0
	g.ufoTimer = 0
	g.bombTimer = 0
	g.state = StatePlaying
}

func (g *Game) update(dt float64) {
	g.updatePlayer(dt)
	g.updateAliens(dt)
	g.updateBullets(dt)
	g.updateUFO(dt)
	g.checkCollisions()
	g.pruneDeadBullets()
	g.checkWinLoseConditions()
}

func (g *Game) updatePlayer(dt float64) {
	if g.player.Invincible {
		g.player.InvTimer -= dt
		if g.player.InvTimer <= 0 {
			g.player.Invincible = false
			g.player.InvTimer = 0
		}
	}
}

func (g *Game) updateAliens(dt float64) {
	alive := g.countAlive()
	if alive == 0 {
		return
	}
	g.alienAnimT += dt
	if g.alienAnimT >= alienAnimInterval {
		g.alienAnimT = 0
		for _, a := range g.aliens {
			if a.Alive {
				a.Frame ^= 1
			}
		}
	}
	speed := alienBaseSpeed(g.level, alive)
	stepInterval := 1.0 / (speed / (AlienWidth + AlienHPadding))
	g.alienTimer += dt
	if g.alienTimer >= stepInterval {
		g.alienTimer = 0
		g.stepAliens()
	}
	g.bombTimer += dt
	interval := bombInterval - float64(g.level-1)*0.1
	if interval < 0.5 {
		interval = 0.5
	}
	if g.bombTimer >= interval {
		g.bombTimer = 0
		g.dropAlienBomb()
	}
}

func (g *Game) stepAliens() {
	minX, maxX := 9999.0, -9999.0
	for _, a := range g.aliens {
		if !a.Alive {
			continue
		}
		if a.Pos.X < minX {
			minX = a.Pos.X
		}
		if a.Pos.X+AlienWidth > maxX {
			maxX = a.Pos.X + AlienWidth
		}
	}
	step := (AlienWidth + AlienHPadding) * 0.25
	if maxX+step*g.alienDirX >= CanvasWidth-10 {
		g.alienDirX = -1
		for _, a := range g.aliens {
			if a.Alive {
				a.Pos.Y += alienDropAmount
			}
		}
	} else if minX+step*g.alienDirX <= 10 {
		g.alienDirX = 1
		for _, a := range g.aliens {
			if a.Alive {
				a.Pos.Y += alienDropAmount
			}
		}
	} else {
		for _, a := range g.aliens {
			if a.Alive {
				a.Pos.X += step * g.alienDirX
			}
		}
	}
}

func (g *Game) dropAlienBomb() {
	bottomAliens := make(map[int]*Alien)
	for _, a := range g.aliens {
		if !a.Alive {
			continue
		}
		prev, ok := bottomAliens[a.Col]
		if !ok || a.Row < prev.Row {
			bottomAliens[a.Col] = a
		}
	}
	if len(bottomAliens) == 0 {
		return
	}
	keys := make([]int, 0, len(bottomAliens))
	for k := range bottomAliens {
		keys = append(keys, k)
	}
	shooter := bottomAliens[keys[g.rng.Intn(len(keys))]]
	g.bullets = append(g.bullets, &Bullet{
		Pos: vec2.Vec2{
			X: shooter.Pos.X + AlienWidth/2 - BulletWidth/2,
			Y: shooter.Pos.Y + AlienHeight,
		},
		Owner: OwnerAlien,
		Alive: true,
	})
}

func (g *Game) updateBullets(dt float64) {
	for _, b := range g.bullets {
		if !b.Alive {
			continue
		}
		if b.Owner == OwnerPlayer {
			b.Pos.Y -= BulletSpeed * dt
			if b.Pos.Y+BulletHeight < 0 {
				b.Alive = false
			}
		} else {
			b.Pos.Y += AlienBombSpeed * dt
			if b.Pos.Y > CanvasHeight {
				b.Alive = false
			}
		}
	}
}

func (g *Game) updateUFO(dt float64) {
	if g.ufo.Active {
		g.ufo.Pos.X += UFOSpeed * g.ufo.DirX * dt
		if g.ufo.Pos.X > CanvasWidth+UFOWidth || g.ufo.Pos.X < -UFOWidth {
			g.ufo.Active = false
			g.ufoTimer = 0
		}
		return
	}
	g.ufoTimer += dt
	if g.ufoTimer >= ufoSpawnInterval {
		g.ufoTimer = 0
		dir := float64(1 - 2*g.rng.Intn(2))
		startX := -UFOWidth
		if dir < 0 {
			startX = CanvasWidth
		}
		g.ufo = UFO{Pos: vec2.Vec2{X: startX, Y: UFOY}, Active: true, DirX: dir}
	}
}

func (g *Game) tryPlayerShoot() {
	count := 0
	for _, b := range g.bullets {
		if b.Alive && b.Owner == OwnerPlayer {
			count++
		}
	}
	if count >= maxPlayerBullets {
		return
	}
	g.bullets = append(g.bullets, &Bullet{
		Pos:   vec2.Vec2{X: g.player.Pos.X + PlayerWidth/2 - BulletWidth/2, Y: g.player.Pos.Y - BulletHeight},
		Owner: OwnerPlayer,
		Alive: true,
	})
}

func (g *Game) checkCollisions() {
	checkBulletAlienCollisions(g.bullets, g.aliens, &g.score,
		func(id int, x, y float64, pts int) {
			if g.score > g.hiScore {
				g.hiScore = g.score
			}
			g.emit(protocol.GameEventMessage{Type: protocol.TypeEvent, Event: protocol.EventAlienKilled, Data: protocol.EventData{X: x, Y: y, Score: pts}})
		},
	)
	checkBulletShieldCollisions(g.bullets, g.shields,
		func(idx int, x, y float64) {
			g.emit(protocol.GameEventMessage{Type: protocol.TypeEvent, Event: protocol.EventShieldHit, Data: protocol.EventData{X: x, Y: y}})
		},
	)
	checkBulletUFOCollision(g.bullets, &g.ufo,
		func(x, y float64, pts int) {
			g.score += pts
			if g.score > g.hiScore {
				g.hiScore = g.score
			}
			g.emit(protocol.GameEventMessage{Type: protocol.TypeEvent, Event: protocol.EventUFOKilled, Data: protocol.EventData{X: x, Y: y, Score: pts}})
		},
		func() int { return ScoreUFOMin + g.rng.Intn((ScoreUFOMax-ScoreUFOMin)/50+1)*50 },
	)
	if checkBulletPlayerCollisions(g.bullets, &g.player) {
		g.player.Lives--
		g.emit(protocol.GameEventMessage{Type: protocol.TypeEvent, Event: protocol.EventPlayerHit, Data: protocol.EventData{X: g.player.Pos.X + PlayerWidth/2, Y: g.player.Pos.Y + PlayerHeight/2}})
		if g.player.Lives <= 0 {
			g.state = StateGameOver
			g.emit(protocol.GameEventMessage{Type: protocol.TypeEvent, Event: protocol.EventGameOver})
		} else {
			g.player.Invincible = true
			g.player.InvTimer = InvincibilityDuration
		}
	}
}

func (g *Game) checkWinLoseConditions() {
	if g.state != StatePlaying {
		return
	}
	if g.countAlive() == 0 {
		g.level++
		if g.level > 10 {
			g.state = StateVictory
			g.emit(protocol.GameEventMessage{Type: protocol.TypeEvent, Event: protocol.EventVictory})
			return
		}
		g.emit(protocol.GameEventMessage{Type: protocol.TypeEvent, Event: protocol.EventLevelComplete, Data: protocol.EventData{Score: g.level}})
		lives := g.player.Lives
		g.startLevel()
		g.player.Lives = lives
	}
	if checkAlienReachedPlayer(g.aliens, g.player.Pos.Y) {
		g.state = StateGameOver
		g.emit(protocol.GameEventMessage{Type: protocol.TypeEvent, Event: protocol.EventGameOver})
	}
}

func (g *Game) countAlive() int {
	n := 0
	for _, a := range g.aliens {
		if a.Alive {
			n++
		}
	}
	return n
}

func (g *Game) pruneDeadBullets() {
	live := g.bullets[:0]
	for _, b := range g.bullets {
		if b.Alive {
			live = append(live, b)
		}
	}
	g.bullets = live
}

func (g *Game) emit(ev protocol.GameEventMessage) {
	if g.onEvent != nil {
		g.onEvent(ev)
	}
}
