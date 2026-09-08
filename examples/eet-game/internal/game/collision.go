package game

// Overlaps returns true when two AABBs intersect.
func (a AABB) Overlaps(b AABB) bool {
	return a.X < b.X+b.W &&
		a.X+a.W > b.X &&
		a.Y < b.Y+b.H &&
		a.Y+a.H > b.Y
}

// checkBulletAlienCollisions tests every live bullet against every live alien.
// Killed aliens are marked dead and the score is incremented. The event
// callback is called for each kill: f(alienID, x, y, score).
func checkBulletAlienCollisions(
	bullets []*Bullet,
	aliens []*Alien,
	score *int,
	onKill func(alienID int, x, y float64, pts int),
) {
	for _, b := range bullets {
		if !b.Alive || b.Owner != OwnerPlayer {
			continue
		}
		bb := b.Bounds()
		for _, a := range aliens {
			if !a.Alive {
				continue
			}
			if bb.Overlaps(a.Bounds()) {
				b.Alive = false
				a.Alive = false
				pts := a.ScoreValue()
				*score += pts
				onKill(a.ID, a.Pos.X+AlienWidth/2, a.Pos.Y+AlienHeight/2, pts)
				break
			}
		}
	}
}

// checkBulletPlayerCollisions tests alien bombs against the player.
// Returns true if the player was hit.
func checkBulletPlayerCollisions(bullets []*Bullet, player *Player) bool {
	if player.Invincible {
		return false
	}
	pb := player.Bounds()
	for _, b := range bullets {
		if !b.Alive || b.Owner != OwnerAlien {
			continue
		}
		if b.Bounds().Overlaps(pb) {
			b.Alive = false
			return true
		}
	}
	return false
}

// checkBulletShieldCollisions tests all bullets against all shields.
// onHit is called with shield index and hit position when a hit occurs.
func checkBulletShieldCollisions(
	bullets []*Bullet,
	shields []*Shield,
	onHit func(idx int, x, y float64),
) {
	for _, b := range bullets {
		if !b.Alive {
			continue
		}
		bb := b.Bounds()
		for i, s := range shields {
			if s.HP <= 0 {
				continue
			}
			if bb.Overlaps(s.Bounds()) {
				b.Alive = false
				s.HP--
				onHit(i, b.Pos.X, b.Pos.Y)
				break
			}
		}
	}
}

// checkBulletUFOCollision tests player bullets against the UFO.
// Returns true and calls onKill if the UFO was hit.
func checkBulletUFOCollision(
	bullets []*Bullet,
	ufo *UFO,
	onKill func(x, y float64, pts int),
	scoreFn func() int,
) {
	if !ufo.Active {
		return
	}
	ub := ufo.Bounds()
	for _, b := range bullets {
		if !b.Alive || b.Owner != OwnerPlayer {
			continue
		}
		if b.Bounds().Overlaps(ub) {
			b.Alive = false
			ufo.Active = false
			pts := scoreFn()
			onKill(ufo.Pos.X+UFOWidth/2, ufo.Pos.Y+UFOHeight/2, pts)
			return
		}
	}
}

// checkAlienReachedPlayer returns true if any live alien has descended below
// the player's Y position (game over by invasion).
func checkAlienReachedPlayer(aliens []*Alien, playerY float64) bool {
	for _, a := range aliens {
		if a.Alive && a.Pos.Y+AlienHeight >= playerY {
			return true
		}
	}
	return false
}
