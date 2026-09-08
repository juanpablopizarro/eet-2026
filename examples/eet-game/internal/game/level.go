package game

import "eet-game/pkg/vec2"

// buildAlienFormation creates the initial alien grid for the given level.
// Higher levels start the formation lower and slightly faster (handled in game.go).
func buildAlienFormation(level int) []*Alien {
	startX := 80.0
	startY := 80.0 + float64(level-1)*4 // descend a bit each level

	aliens := make([]*Alien, 0, AlienRows*AlienCols)
	id := 0
	for row := 0; row < AlienRows; row++ {
		var aType AlienType
		switch {
		case row == AlienRows-1:
			aType = AlienTypeTop
		case row >= AlienRows-3:
			aType = AlienTypeMiddle
		default:
			aType = AlienTypeBottom
		}

		for col := 0; col < AlienCols; col++ {
			x := startX + float64(col)*(AlienWidth+AlienHPadding)
			y := startY + float64(row)*(AlienHeight+AlienVPadding)
			aliens = append(aliens, &Alien{
				ID:    id,
				Pos:   vec2.Vec2{X: x, Y: y},
				Row:   row,
				Col:   col,
				Type:  aType,
				Alive: true,
				Frame: 0,
			})
			id++
		}
	}
	return aliens
}

// buildShields creates the four protective barriers.
func buildShields() []*Shield {
	ys := CanvasHeight - 120.0
	positions := []float64{120, 270, 450, 600}
	shields := make([]*Shield, len(positions))
	for i, x := range positions {
		shields[i] = &Shield{
			Pos: vec2.Vec2{X: x, Y: ys},
			HP:  ShieldMaxHP,
		}
	}
	return shields
}

// alienBaseSpeed returns the horizontal speed (px/s) for the alien formation
// at the given level and with the given number of remaining aliens.
func alienBaseSpeed(level int, remaining int) float64 {
	base := 28.0 + float64(level-1)*8
	// Formation speeds up as aliens are eliminated
	ratio := 1.0 - float64(remaining)/float64(AlienRows*AlienCols)
	return base + ratio*120
}
