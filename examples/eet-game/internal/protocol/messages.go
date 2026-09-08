// Package protocol defines the WebSocket message types exchanged between
// the game server and the browser client.
package protocol

// MessageType is the discriminator field in every message.
type MessageType string

const (
	// Client → Server
	TypeInput MessageType = "input"

	// Server → Client
	TypeState MessageType = "state"
	TypeEvent MessageType = "event"
)

// InputAction represents a player action sent from the client.
type InputAction string

const (
	ActionMoveLeft  InputAction = "move_left"
	ActionMoveRight InputAction = "move_right"
	ActionShoot     InputAction = "shoot"
	ActionStart     InputAction = "start"
	ActionPause     InputAction = "pause"
	ActionRestart   InputAction = "restart"
)

// InputMessage is sent by the client to communicate player input.
type InputMessage struct {
	Type   MessageType `json:"type"`
	Action InputAction `json:"action"`
}

// PlayerState holds the serialisable state of the player's ship.
type PlayerState struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Lives int     `json:"lives"`
	// Invincible is true during the brief grace period after being hit.
	Invincible bool `json:"invincible"`
}

// AlienState holds the serialisable state of a single alien.
type AlienState struct {
	ID    int     `json:"id"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Type  int     `json:"alienType"` // 1=bottom, 2=middle, 3=top
	Alive bool    `json:"alive"`
	Frame int     `json:"frame"` // animation frame 0 or 1
}

// BulletState holds the serialisable state of a bullet / bomb.
type BulletState struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Owner string  `json:"owner"` // "player" or "alien"
}

// ShieldState holds the serialisable state of a defensive shield.
type ShieldState struct {
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	HP int     `json:"hp"` // 0 = destroyed
}

// UFOState holds the serialisable state of the mystery UFO.
type UFOState struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Active bool    `json:"active"`
}

// GameStateMessage is sent from the server every tick containing the full
// authoritative game state.
type GameStateMessage struct {
	Type    MessageType   `json:"type"`
	Player  PlayerState   `json:"player"`
	Aliens  []AlienState  `json:"aliens"`
	Bullets []BulletState `json:"bullets"`
	Shields []ShieldState `json:"shields"`
	UFO     UFOState      `json:"ufo"`
	Score   int           `json:"score"`
	HiScore int           `json:"hiScore"`
	Level   int           `json:"level"`
	State   string        `json:"state"` // "menu"|"playing"|"paused"|"game_over"|"victory"
}

// GameEventKind is the category of a server-pushed event.
type GameEventKind string

const (
	EventAlienKilled    GameEventKind = "alien_killed"
	EventPlayerHit      GameEventKind = "player_hit"
	EventUFOKilled      GameEventKind = "ufo_killed"
	EventLevelComplete  GameEventKind = "level_complete"
	EventGameOver       GameEventKind = "game_over"
	EventVictory        GameEventKind = "victory"
	EventShieldHit      GameEventKind = "shield_hit"
)

// EventData carries extra information about a game event.
type EventData struct {
	X     float64 `json:"x,omitempty"`
	Y     float64 `json:"y,omitempty"`
	Score int     `json:"score,omitempty"`
}

// GameEventMessage is sent from the server whenever a notable game event occurs.
type GameEventMessage struct {
	Type  MessageType   `json:"type"`
	Event GameEventKind `json:"event"`
	Data  EventData     `json:"data"`
}
