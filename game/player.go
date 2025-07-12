package game

import (
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type Player struct {
	Score           int
	Level           int
	CurrentPolymino *Polyomino
	NextPolyomino   *Polyomino
	LinesCleared    int
	HasSwapped      bool // Track if player has already swapped the current block
}

func NewPlayer() *Player {
	nextPolyomino := GeneratePolyomino()

	return &Player{
		Score:           0,
		Level:           1,
		CurrentPolymino: nil,
		NextPolyomino:   nextPolyomino,
		LinesCleared:    0,
		HasSwapped:      false,
	}
}

func GeneratePolyomino() *Polyomino {
	blockPositions := []Position{}
	size := rng.Intn(5) + 1
	blockPositions = append(blockPositions, Position{X: 0, Y: 0})
	for i := 0; i < size; i++ {
		potentialPositionsSize := len(generateBlockOptions(blockPositions))
		if potentialPositionsSize == 0 {
			break
		}
		blockPos := rng.Intn(potentialPositionsSize)
		blockPositions = append(blockPositions, generateBlockOptions(blockPositions)[blockPos])
	}

	colors := []string{
		"blue",
		"red",
		"green",
		"yellow",
		"cyan",
		"magenta",
		"white",
	}
	color := colors[rng.Intn(len(colors))]

	blocks := make([]Block, len(blockPositions))
	for i, pos := range blockPositions {
		blocks[i] = Block{
			Position: pos,
			Color:    color,
		}
	}

	lowestPosition := GetLowestBlockPosition(blocks)

	return &Polyomino{
		Blocks:   blocks,
		Position: Position{X: GameFieldWidth / 2, Y: -lowestPosition.Y - 1},
		Placed:   false,
	}
}

func GetLowestBlockPosition(blocks []Block) Position {
	if len(blocks) == 0 {
		return Position{X: 0, Y: 0}
	}

	lowest := blocks[0].Position
	for _, block := range blocks {
		if block.Position.Y > lowest.Y {
			lowest = block.Position
		}
	}

	return lowest
}

func generateBlockOptions(currentPositions []Position) []Position {
	potentialPositions := make(map[Position]bool)

	for _, position := range currentPositions {
		adjacentPositions := []Position{
			{X: position.X + 1, Y: position.Y},
			{X: position.X - 1, Y: position.Y},
			{X: position.X, Y: position.Y + 1},
			{X: position.X, Y: position.Y - 1},
		}

		for _, adj := range adjacentPositions {
			potentialPositions[adj] = true
		}
	}

	for _, position := range currentPositions {
		delete(potentialPositions, position)
	}

	result := make([]Position, 0, len(potentialPositions))
	for pos := range potentialPositions {
		result = append(result, pos)
	}

	return result

}
