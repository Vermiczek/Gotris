package game

import "fmt"

func (g *Game) SwapBlocks() {
	if g.Player.CurrentPolymino == nil || g.Player.HasSwapped {
		GetLoggerInstance().Log("Cannot swap - already swapped or no active block")
		return
	}

	originalPiece := g.Player.CurrentPolymino
	originalNextPiece := g.Player.NextPolyomino

	g.Player.CurrentPolymino = originalNextPiece
	g.Player.NextPolyomino = originalPiece

	lowestY := GetLowestBlockPosition(g.Player.CurrentPolymino.Blocks).Y
	g.Player.CurrentPolymino.Position = Position{
		X: GameFieldWidth / 2,
		Y: -lowestY - 1,
	}

	lowestNextY := GetLowestBlockPosition(g.Player.NextPolyomino.Blocks).Y
	g.Player.NextPolyomino.Position = Position{
		X: GameFieldWidth / 2,
		Y: -lowestNextY - 1,
	}

	if g.checkCollision() {
		g.Player.CurrentPolymino = originalPiece
		g.Player.NextPolyomino = originalNextPiece

		GetLoggerInstance().Log("Cannot swap - collision detected")
	} else {
		g.Player.HasSwapped = true
		GetLoggerInstance().Log("Blocks swapped successfully")
	}
}

func (g *Game) HardDrop() {
	if g.Player.CurrentPolymino == nil {
		return
	}

	movesMade := 0
	for !g.checkMovementCollision(0, 1) {
		g.Player.CurrentPolymino.Move(0, 1)
		movesMade++
	}

	g.placeCurrentPolyomino()

	GetLoggerInstance().Log("Hard dropped block by " + fmt.Sprintf("%d", movesMade) + " rows")
}
