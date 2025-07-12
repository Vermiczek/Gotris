package game

import (
	"time"
)

const (
	GameFieldStartX = 3
	GameFieldWidth  = 15
	GameFieldEndX   = GameFieldStartX + GameFieldWidth

	GameFieldStartY = 1
	GameFieldHeight = 25
	GameFieldEndY   = GameFieldStartY + GameFieldHeight

	BlockWidth = 2
)

type Game struct {
	timer        *GameTimer
	eventHandler *EventHandler
	Player       *Player
	placedBlocks []Block
	lastDropTime int64
	Scoring      *ScoringSystem
	UI           *Interface
	IsGameOver   bool // Flag to indicate if the game is over
}

func NewGame() *Game {
	renderer := GetRendererInstance()

	return &Game{
		timer:        NewGameTimer(),
		eventHandler: NewEventHandler(),
		Player:       NewPlayer(),
		Scoring:      NewScoringSystem(),
		UI:           NewInterface(renderer),
	}
}

func (g *Game) Start() {
	renderer := GetRendererInstance()
	g.eventHandler.Start()
	g.timer.Reset()

	defer g.eventHandler.Stop()

	running := true
	for running {
		g.Update()

		if g.IsGameOver {
			select {
			case <-g.eventHandler.QuitChannel():
				running = false
			case event := <-g.eventHandler.InputEvents:
				if event.Action == "restart" {
					g.Reset()
				} else if event.Action == "quit" {
					running = false
				}
			default:
				renderer.Render()
				time.Sleep(100 * time.Millisecond)
			}
		} else {
			// Normal game loop
			select {
			case <-g.eventHandler.QuitChannel():
				running = false
			case event := <-g.eventHandler.InputEvents:
				g.processInput(event)
			default:
				renderer.Render()
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func (g *Game) Update() {
	rendererInstance.RenderGame(g)
	g.timer.Update()

	currentTime := g.timer.elapsed

	g.drop(currentTime)
}

func (g *Game) drop(currentTime int64) {
	dropInterval := g.Scoring.GetDropSpeed()

	if g.Player.CurrentPolymino != nil {
		if currentTime-g.lastDropTime >= int64(dropInterval) {
			canMoveDown := true

			for _, block := range g.Player.CurrentPolymino.Blocks {
				absX := g.Player.CurrentPolymino.Position.X + block.Position.X
				absY := g.Player.CurrentPolymino.Position.Y + block.Position.Y

				isFloor := absY+1 >= GameFieldHeight
				isTouchingPlacedBlock := false

				for _, placedBlock := range g.placedBlocks {
					if placedBlock.Position.X == absX && placedBlock.Position.Y == absY+1 {
						isTouchingPlacedBlock = true
						break
					}
				}

				if isFloor || isTouchingPlacedBlock {
					canMoveDown = false
					break
				}
			}

			if canMoveDown {
				g.Player.CurrentPolymino.Move(0, 1)
			} else {
				g.placeCurrentPolyomino()
			}

			g.lastDropTime = currentTime
		}
	} else {
		g.Player.CurrentPolymino = g.Player.NextPolyomino

		g.Player.NextPolyomino = GeneratePolyomino()

		g.lastDropTime = currentTime
	}
}

func (g *Game) placeCurrentPolyomino() {
	if g.Player.CurrentPolymino == nil {
		return
	}

	for _, block := range g.Player.CurrentPolymino.Blocks {
		absY := g.Player.CurrentPolymino.Position.Y + block.Position.Y
		if absY < 0 {
			GetLoggerInstance().Log("GAME OVER!")
			g.IsGameOver = true
		}
	}

	g.Player.HasSwapped = false

	for _, block := range g.Player.CurrentPolymino.Blocks {
		absY := g.Player.CurrentPolymino.Position.Y + block.Position.Y

		if absY >= 0 {
			placedBlock := Block{
				Position: Position{
					X: g.Player.CurrentPolymino.Position.X + block.Position.X,
					Y: absY,
				},
				Color: block.Color,
			}

			g.placedBlocks = append(g.placedBlocks, placedBlock)
		}
	}

	g.CheckLineClear()

	g.Player.CurrentPolymino = nil
}

func (g *Game) TryRotate() {
	if g.Player.CurrentPolymino == nil {
		return
	}

	originalBlocks := make([]Block, len(g.Player.CurrentPolymino.Blocks))
	for i, block := range g.Player.CurrentPolymino.Blocks {
		originalBlocks[i] = Block{
			Position: Position{X: block.Position.X, Y: block.Position.Y},
			Color:    block.Color,
		}
	}
	originalX := g.Player.CurrentPolymino.Position.X
	originalY := g.Player.CurrentPolymino.Position.Y

	g.Player.CurrentPolymino.Rotate(true)

	if !g.checkCollision() {
		return
	}

	_, originalCollisionType := g.checkCollisionWithType()

	if originalCollisionType == "block" {
		g.Player.CurrentPolymino.Position.X = originalX
		g.Player.CurrentPolymino.Position.Y = originalY
		for i := range g.Player.CurrentPolymino.Blocks {
			g.Player.CurrentPolymino.Blocks[i].Position.X = originalBlocks[i].Position.X
			g.Player.CurrentPolymino.Blocks[i].Position.Y = originalBlocks[i].Position.Y
		}
		return
	}

	kickOffsets := []struct{ x, y int }{
		{1, 0},
		{-1, 0},
		{0, -1},
		{2, 0},
		{-2, 0},
	}

	for _, offset := range kickOffsets {
		g.Player.CurrentPolymino.Position.X = originalX + offset.x
		g.Player.CurrentPolymino.Position.Y = originalY + offset.y

		collision, _ := g.checkCollisionWithType()

		if !collision {
			return
		}
	}

	g.Player.CurrentPolymino.Position.X = originalX
	g.Player.CurrentPolymino.Position.Y = originalY

	for i := range g.Player.CurrentPolymino.Blocks {
		g.Player.CurrentPolymino.Blocks[i].Position.X = originalBlocks[i].Position.X
		g.Player.CurrentPolymino.Blocks[i].Position.Y = originalBlocks[i].Position.Y
	}
}

func (g *Game) checkCollisionWithType() (bool, string) {
	if g.Player.CurrentPolymino == nil {
		return false, ""
	}

	for _, block := range g.Player.CurrentPolymino.Blocks {
		absX := g.Player.CurrentPolymino.Position.X + block.Position.X
		absY := g.Player.CurrentPolymino.Position.Y + block.Position.Y

		if absX < 0 || absX >= GameFieldWidth || absY >= GameFieldHeight {
			return true, "wall"
		}

		for _, placedBlock := range g.placedBlocks {
			if absY >= 0 && placedBlock.Position.X == absX && placedBlock.Position.Y == absY {
				return true, "block"
			}
		}
	}

	return false, ""
}

func (g *Game) checkCollision() bool {
	collision, _ := g.checkCollisionWithType()
	return collision
}

func (g *Game) checkMovementCollision(dx, dy int) bool {
	if g.Player.CurrentPolymino == nil {
		return false
	}

	for _, block := range g.Player.CurrentPolymino.Blocks {
		newX := g.Player.CurrentPolymino.Position.X + block.Position.X + dx
		newY := g.Player.CurrentPolymino.Position.Y + block.Position.Y + dy

		if newX < 0 || newX >= GameFieldWidth || newY >= GameFieldHeight {
			return true
		}

		for _, placedBlock := range g.placedBlocks {
			if newY >= 0 && placedBlock.Position.X == newX && placedBlock.Position.Y == newY {
				return true
			}
		}
	}

	return false
}

func (g *Game) processInput(event Event) {
	logger.Log("Processing input event: " + event.Action)

	if g.Player.CurrentPolymino == nil {
		return
	}

	switch event.Action {
	case "up":
		g.TryRotate()
	case "down":
		if !g.checkMovementCollision(0, 1) {
			g.Player.CurrentPolymino.Move(0, 1)
		}
	case "left":
		if !g.checkMovementCollision(-1, 0) {
			g.Player.CurrentPolymino.Move(-1, 0)
		}
	case "right":
		if !g.checkMovementCollision(1, 0) {
			g.Player.CurrentPolymino.Move(1, 0)
		}
	case "space":
		g.TryRotate()
	case "swap":
		g.SwapBlocks()
	case "hardDrop":
		g.HardDrop()
	}
}
