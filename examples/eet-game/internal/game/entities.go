// Package game contains the core game logic for the Space Invaders clone.
package game

import "eet-game/pkg/vec2"

// ---- Constants -----------------------------------------------------------

const (
	CanvasWidth  = 800
	CanvasHeight = 600

	PlayerWidth  = 40.0
	PlayerHeight = 24.0
	PlayerSpeed  = 220.0 // pixels per second

	BulletWidth   = 4.0
	BulletHeight  = 12.0
	BulletSpeed   = 380.0 // pixels per second (player bullet)
	AlienBombSpeed = 160.0 // pixels per second (alien bomb)

	AlienWidth  = 36.0
	AlienHeight = 28.0
	AlienCols   = 11
	AlienRows   = 5

	AlienHPadding = 16.0
	AlienVPadding = 14.0

	ShieldWidth  = 64.0
	ShieldHeight = 40.0
	ShieldMaxHP  = 4

	UFOY      = 40.0
	UFOWidth  = 48.0
	UFOHeight = 22.0
	UFOSpeed  = 100.0

	InvincibilityDuration = 2.0 // seconds after being hit

	// Score values per alien row (row 0 = bottom)
	ScoreRow0 = 10
	ScoreRow1 = 10
	ScoreRow2 = 20
	ScoreRow3 = 20
	ScoreRow4 = 30 // top row

	ScoreUFOMin = 50
	ScoreUFOMax = 300

	alienDropAmount = 20.0 // pixels aliens drop when reversing direction
)

// ---- Player --------------------------------------------------------------

// Player represents the human-controlled spaceship.
type Player struct {
	Pos        vec2.Vec2
	Lives      int
	Invincible bool
	InvTimer   float64 // counts down to 0
}

// NewPlayer creates a freshly spawned player centred at the bottom of the canvas.
func NewPlayer() Player {
	return Player{
		Pos: vec2.Vec2{
			X: CanvasWidth/2 - PlayerWidth/2,
			Y: CanvasHeight - 60,
		},
		Lives: 3,
	}
}

// Bounds returns the AABB of the player.
func (p *Player) Bounds() AABB {
	return AABB{X: p.Pos.X, Y: p.Pos.Y, W: PlayerWidth, H: PlayerHeight}
}

// ---- Alien ---------------------------------------------------------------

// AlienType determines the point value and sprite used for an alien.
type AlienType int

const (
	AlienTypeBottom AlienType = 1 // row 0-1 (10 pts)
	AlienTypeMiddle AlienType = 2 // row 2-3 (20 pts)
	AlienTypeTop    AlienType = 3 // row 4   (30 pts)
)

// Alien is a single enemy in the formation.
type Alien struct {
	ID    int
	Pos   vec2.Vec2
	Row   int
	Col   int
	Type  AlienType
	Alive bool
	Frame int // 0 or 1 – toggles for walk animation
}

// Bounds returns the AABB for the alien.
func (a *Alien) Bounds() AABB {
	return AABB{X: a.Pos.X, Y: a.Pos.Y, W: AlienWidth, H: AlienHeight}
}

// ScoreValue returns the point value for destroying this alien.
func (a *Alien) ScoreValue() int {
	switch a.Type {
	case AlienTypeTop:
		return ScoreRow4
	case AlienTypeMiddle:
		return ScoreRow2
	default:
		return ScoreRow0
	}
}

// ---- Bullet --------------------------------------------------------------

// BulletOwner identifies who fired a projectile.
type BulletOwner string

const (
	OwnerPlayer BulletOwner = "player"
	OwnerAlien  BulletOwner = "alien"
)

// Bullet represents a projectile in flight.
type Bullet struct {
	Pos   vec2.Vec2
	Owner BulletOwner
	Alive bool
}

// Bounds returns the AABB for the bullet.
func (b *Bullet) Bounds() AABB {
	return AABB{X: b.Pos.X, Y: b.Pos.Y, W: BulletWidth, H: BulletHeight}
}

// ---- Shield --------------------------------------------------------------

// Shield is a destructible barrier that the player can hide behind.
type Shield struct {
	Pos vec2.Vec2
	HP  int
}

// Bounds returns the AABB for the shield.
func (s *Shield) Bounds() AABB {
	return AABB{X: s.Pos.X, Y: s.Pos.Y, W: ShieldWidth, H: ShieldHeight}
}

// ---- UFO -----------------------------------------------------------------

// UFO is the mystery bonus saucer that flies across the top of the screen.
type UFO struct {
	Pos    vec2.Vec2
	Active bool
	DirX   float64 // +1 or -1
}

// Bounds returns the AABB for the UFO.
func (u *UFO) Bounds() AABB {
	return AABB{X: u.Pos.X, Y: u.Pos.Y, W: UFOWidth, H: UFOHeight}
}

// ---- AABB ----------------------------------------------------------------

// AABB is an axis-aligned bounding box used for collision detection.
type AABB struct {
	X, Y, W, H float64
}
