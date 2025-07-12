package game

import "math"

type Block struct {
	Position
	Color string
}

type Position struct {
	X int
	Y int
}

func NewBlock(x, y int, color string) *Block {
	return &Block{
		Position: Position{
			X: x,
			Y: y,
		},
		Color: color,
	}
}

type Polyomino struct {
	Blocks []Block
	Position
	Placed        bool
	RotationPoint Position
}

func NewPolyomino(blocks []Block, x, y int, placed bool) *Polyomino {
	return &Polyomino{
		Blocks: blocks,
		Position: Position{
			X: x,
			Y: y,
		},
		Placed:        placed,
		RotationPoint: CalculateRotationPoint(blocks),
	}
}

func CalculateRotationPoint(blocks []Block) Position {
	var sumX, sumY float64
	for _, block := range blocks {
		sumX += float64(block.Position.X)
		sumY += float64(block.Position.Y)
	}

	centerX := sumX / float64(len(blocks))
	centerY := sumY / float64(len(blocks))

	return Position{
		X: int(math.Round(centerX)),
		Y: int(math.Round(centerY)),
	}
}

func (t *Polyomino) Move(dx, dy int) {
	t.Position.X += dx
	t.Position.Y += dy
}

func (p *Polyomino) Rotate(clockwise bool) {
	localPivotX := float64(p.RotationPoint.X)
	localPivotY := float64(p.RotationPoint.Y)

	for i := 0; i < len(p.Blocks); i++ {
		relX := float64(p.Blocks[i].Position.X) - localPivotX
		relY := float64(p.Blocks[i].Position.Y) - localPivotY

		var newX, newY float64

		if clockwise {
			newX = localPivotX + relY
			newY = localPivotY - relX
		} else {
			newX = localPivotX - relY
			newY = localPivotY + relX
		}

		p.Blocks[i].Position.X = int(math.Round(newX))
		p.Blocks[i].Position.Y = int(math.Round(newY))
	}
}
