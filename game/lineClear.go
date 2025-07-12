package game

import "fmt"

func (g *Game) CheckLineClear() {
	clearedLines := 0

	blocksInRow := make(map[int]int)

	for _, block := range g.placedBlocks {
		blocksInRow[block.Position.Y]++
	}

	fullRows := []int{}
	for y, count := range blocksInRow {
		if count >= GameFieldWidth {
			fullRows = append(fullRows, y)
			clearedLines++
		}
	}

	if clearedLines == 0 {
		return
	}

	newPlacedBlocks := []Block{}
	for _, block := range g.placedBlocks {
		isInFullRow := false
		for _, fullY := range fullRows {
			if block.Position.Y == fullY {
				isInFullRow = true
				break
			}
		}

		if !isInFullRow {
			newPlacedBlocks = append(newPlacedBlocks, block)
		}
	}

	for _, fullY := range fullRows {
		for i := range newPlacedBlocks {
			if newPlacedBlocks[i].Position.Y < fullY {
				newPlacedBlocks[i].Position.Y++
			}
		}
	}

	g.placedBlocks = newPlacedBlocks

	droppedRows := g.DropPlacedBlocks()

	g.Scoring.AddLines(clearedLines)

	GetLoggerInstance().Log(fmt.Sprintf("Cleared %d lines! Dropped blocks by %d rows. Score: %d, Level: %d",
		clearedLines, droppedRows, g.Scoring.Score, g.Scoring.Level))
}

func (g *Game) DropPlacedBlocks() int {
	if len(g.placedBlocks) == 0 {
		return 0
	}

	movesMade := 0
	somethingMoved := true

	for somethingMoved {
		somethingMoved = false

		occupiedSpaces := make(map[string]bool)
		for _, block := range g.placedBlocks {
			key := fmt.Sprintf("%d,%d", block.Position.X, block.Position.Y)
			occupiedSpaces[key] = true
		}

		anyBlockNearBottom := false
		for _, block := range g.placedBlocks {
			if block.Position.Y >= GameFieldHeight-2 {
				anyBlockNearBottom = true
				break
			}
		}

		if anyBlockNearBottom {
			break
		}

		canMove := make([]bool, len(g.placedBlocks))

		for i := range g.placedBlocks {
			blockX := g.placedBlocks[i].Position.X
			blockY := g.placedBlocks[i].Position.Y

			if blockY >= GameFieldHeight-1 {
				canMove[i] = false
				continue
			}

			belowKey := fmt.Sprintf("%d,%d", blockX, blockY+1)
			if occupiedSpaces[belowKey] {
				canMove[i] = false
				continue
			}

			canMove[i] = true
			somethingMoved = true
		}

		if somethingMoved {
			for i, block := range g.placedBlocks {
				if canMove[i] {
					key := fmt.Sprintf("%d,%d", block.Position.X, block.Position.Y)
					delete(occupiedSpaces, key)
				}
			}

			for i := range g.placedBlocks {
				if canMove[i] {
					g.placedBlocks[i].Position.Y++
				}
			}

			movesMade++
		}
	}

	GetLoggerInstance().Log(fmt.Sprintf("Dropped blocks by %d rows", movesMade))
	return movesMade
}
