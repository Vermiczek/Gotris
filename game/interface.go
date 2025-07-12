package game

import "fmt"

type Interface struct {
	renderer   *Renderer
	interfaceX int
	interfaceY int
	width      int
	height     int
	separators []int
}

func NewInterface(r *Renderer) *Interface {
	// Calculate the right wall of the game field as the start of the interface
	gameRightWallScreenX := GameFieldStartX + (GameFieldWidth * BlockWidth)

	return &Interface{
		renderer:   r,
		interfaceX: gameRightWallScreenX + 2,
		interfaceY: 1,
		width:      r.ScreenWidth - gameRightWallScreenX - 3,
		height:     GameFieldHeight,
		separators: []int{2, 5, 8, 11},
	}
}

func (ui *Interface) Draw(game *Game) {
	ui.Clear()
	ui.DrawSeparators()

	ui.DrawTimeSection(game.timer.elapsed / 1000)
	ui.DrawLevelSection(game.Scoring.Level)
	ui.DrawLinesSection(game.Scoring.LinesCleared)
	ui.DrawScoreSection(game.Scoring.Score)
	ui.DrawNextSection(game.Player.NextPolyomino)

	if game.IsGameOver {
		ui.DrawGameOverScreen()
	}
}

func (ui *Interface) Clear() {
	for y := ui.interfaceY; y < ui.interfaceY+ui.height; y++ {
		for x := ui.interfaceX; x < ui.interfaceX+ui.width; x++ {
			if x < ui.renderer.ScreenWidth-1 && y < ui.renderer.ScreenHeight-1 {
				ui.renderer.Pixels[y][x] = ColoredPixel{Char: ' ', Color: ""}
			}
		}
	}
}

func (ui *Interface) DrawSeparators() {
	for _, y := range ui.separators {
		if ui.interfaceY+y < ui.renderer.ScreenHeight-1 {
			ui.DrawSeparator(y)
		}
	}
}

func (ui *Interface) DrawSeparator(y int) {
	// Draw a T-junction on the left side where the separator meets the game border
	ui.renderer.Pixels[ui.interfaceY+y][ui.interfaceX-1] = ColoredPixel{Char: '╠', Color: "cyan"}

	// Draw the separator line
	for x := 0; x < ui.width; x++ {
		ui.renderer.Pixels[ui.interfaceY+y][ui.interfaceX+x] = ColoredPixel{Char: '═', Color: "cyan"}
	}

	// Draw a T-junction on the right side where the separator meets the outer border
	if ui.interfaceX+ui.width < ui.renderer.ScreenWidth-1 {
		ui.renderer.Pixels[ui.interfaceY+y][ui.interfaceX+ui.width] = ColoredPixel{Char: '╣', Color: "cyan"}
	}
}

func (ui *Interface) DrawLabel(text string, yOffset int, color string) {
	for i, char := range text {
		if ui.interfaceX+i < ui.renderer.ScreenWidth-1 {
			ui.renderer.Pixels[ui.interfaceY+yOffset][ui.interfaceX+i] = ColoredPixel{Char: char, Color: color}
		}
	}
}

func (ui *Interface) DrawTimeSection(seconds int64) {
	ui.DrawLabel("TIME", 0, "white")
	ui.DrawLabel(fmt.Sprintf("%d", seconds), 1, "yellow")
}

func (ui *Interface) DrawLevelSection(level int) {
	ui.DrawLabel("LEVEL", 3, "white")
	ui.DrawLabel(fmt.Sprintf("%d", level), 4, "green")
}

func (ui *Interface) DrawLinesSection(lines int) {
	ui.DrawLabel("LINES", 6, "white")
	ui.DrawLabel(fmt.Sprintf("%d", lines), 7, "cyan")
}

func (ui *Interface) DrawScoreSection(score int) {
	ui.DrawLabel("SCORE", 9, "white")
	ui.DrawLabel(fmt.Sprintf("%d", score), 10, "yellow")
}

func (ui *Interface) DrawFieldInfoSection(width, height int) {
	ui.DrawLabel("FIELD", 12, "white")
	ui.DrawLabel(fmt.Sprintf("%dx%d", width, height), 13, "cyan")
}

func (ui *Interface) DrawNextSection(next *Polyomino) {
	displayWidth := 8
	displayHeight := 8

	// Calculate center position
	startX := ui.interfaceX + (ui.width-displayWidth)/2
	startY := ui.interfaceY + ui.separators[len(ui.separators)-1] + 2

	// Draw section title
	for i, c := range "NEXT" {
		ui.renderer.Pixels[startY-1][ui.interfaceX+i] = ColoredPixel{Char: c, Color: "white"}
	}

	// Draw a bordered box for the NEXT block
	// Top border
	ui.renderer.Pixels[startY][startX-1] = ColoredPixel{Char: '╔', Color: "cyan"} // Top-left corner
	for x := 0; x < displayWidth; x++ {
		ui.renderer.Pixels[startY][startX+x] = ColoredPixel{Char: '═', Color: "cyan"}
	}
	ui.renderer.Pixels[startY][startX+displayWidth] = ColoredPixel{Char: '╗', Color: "cyan"} // Top-right corner

	// Side borders
	for y := 1; y < displayHeight; y++ {
		ui.renderer.Pixels[startY+y][startX-1] = ColoredPixel{Char: '║', Color: "cyan"}            // Left border
		ui.renderer.Pixels[startY+y][startX+displayWidth] = ColoredPixel{Char: '║', Color: "cyan"} // Right border
	}

	// Bottom border
	ui.renderer.Pixels[startY+displayHeight][startX-1] = ColoredPixel{Char: '╚', Color: "cyan"} // Bottom-left corner
	for x := 0; x < displayWidth; x++ {
		ui.renderer.Pixels[startY+displayHeight][startX+x] = ColoredPixel{Char: '═', Color: "cyan"}
	}
	ui.renderer.Pixels[startY+displayHeight][startX+displayWidth] = ColoredPixel{Char: '╝', Color: "cyan"} // Bottom-right corner

	// Clear the display area (inside the border)
	for y := 0; y < displayHeight-1; y++ {
		for x := 0; x < displayWidth; x++ {
			ui.renderer.Pixels[startY+y+1][startX+x] = ColoredPixel{Char: ' ', Color: ""}
		}
	}

	if next == nil {
		return
	}

	// Calculate polyomino bounds
	var minX, maxX, minY, maxY int
	for _, block := range next.Blocks {
		if block.Position.X < minX {
			minX = block.Position.X
		}
		if block.Position.X > maxX {
			maxX = block.Position.X
		}
		if block.Position.Y < minY {
			minY = block.Position.Y
		}
		if block.Position.Y > maxY {
			maxY = block.Position.Y
		}
	}

	width := maxX - minX + 1
	height := maxY - minY + 1

	// Calculate centering offsets - adjusted for the border and double-width blocks
	offsetX := (displayWidth - (width * 2)) / 2
	offsetY := (displayHeight - height) / 2

	// Draw the polyomino centered
	for _, block := range next.Blocks {
		baseX := (block.Position.X-minX)*2 + offsetX
		baseY := block.Position.Y - minY + offsetY

		if baseY >= 0 && baseY < displayHeight {
			// First block character
			if baseX >= 0 && baseX < displayWidth {
				ui.renderer.Pixels[startY+baseY+1][startX+baseX] = ColoredPixel{Char: '█', Color: block.Color}
			}

			// Second block character to make it double width
			if baseX+1 >= 0 && baseX+1 < displayWidth {
				ui.renderer.Pixels[startY+baseY+1][startX+baseX+1] = ColoredPixel{Char: '█', Color: block.Color}
			}
		}
	}
}
