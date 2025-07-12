package game

import (
	"fmt"
	"sync"
)

type ColoredPixel struct {
	Char  rune
	Color string
}

type Renderer struct {
	ScreenWidth  int
	ScreenHeight int
	Pixels       [][]ColoredPixel
	Timer        int
}

var rendererInstance *Renderer
var rendererOnce sync.Once

func GetRendererInstance() *Renderer {
	rendererOnce.Do(func() {
		width := GameFieldStartX + (GameFieldWidth * BlockWidth) + 20
		height := GameFieldEndY + 3

		rendererInstance = &Renderer{
			ScreenWidth:  width,
			ScreenHeight: height,
			Pixels:       make([][]ColoredPixel, height),
		}

		for i := range rendererInstance.Pixels {
			rendererInstance.Pixels[i] = make([]ColoredPixel, width)
			for j := range rendererInstance.Pixels[i] {
				rendererInstance.Pixels[i][j] = ColoredPixel{Char: ' ', Color: ""}
			}
		}

		rendererInstance.BuildBorder()
	})
	return rendererInstance
}

func (r *Renderer) BuildBorder() {
	gameRightWallScreenX := GameFieldStartX + (GameFieldWidth * BlockWidth)

	interfaceWidth := 15

	if gameRightWallScreenX+interfaceWidth >= r.ScreenWidth {
		GetLoggerInstance().Log("Warning: Screen width may be too small for proper interface display")
	}

	for y := 0; y < r.ScreenHeight; y++ {
		for x := 0; x < r.ScreenWidth; x++ {
			// Outer screen border
			isOuterBorder := y == 0 || y == r.ScreenHeight-1 || x == 0 || x == r.ScreenWidth-1

			// Game field walls (left wall and right wall)
			isLeftGameWall := x == GameFieldStartX-1 && y > 0 && y < GameFieldEndY
			isRightGameWall := x == gameRightWallScreenX && y > 0 && y < GameFieldEndY

			// Game field floor (bottom wall)
			isGameFloor := y == GameFieldEndY && x > GameFieldStartX-1 && x < gameRightWallScreenX+1

			// Game field corners
			isGameTopLeft := x == GameFieldStartX-1 && y == GameFieldStartY-1
			isGameTopRight := x == gameRightWallScreenX && y == GameFieldStartY-1
			isGameBottomLeft := x == GameFieldStartX-1 && y == GameFieldEndY
			isGameBottomRight := x == gameRightWallScreenX && y == GameFieldEndY

			if isOuterBorder {
				// Customze the outer border appearance
				var borderChar rune

				if x == GameFieldStartX-1 && y == 0 {
					borderChar = '╦' // Connect top border to game's left wall
				} else if x == gameRightWallScreenX && y == 0 {
					borderChar = '╦' // Connect top border to game's right wall
				} else if x == GameFieldStartX-1 && y == r.ScreenHeight-1 {
					borderChar = '╩' // Connect bottom border to game's left wall
				} else if x == gameRightWallScreenX && y == r.ScreenHeight-1 {
					borderChar = '╩' // Connect bottom border to game's right wall
				} else if x == 0 && y == 0 {
					borderChar = '╔' // Top-left corner
				} else if x == r.ScreenWidth-1 && y == 0 {
					borderChar = '╗' // Top-right corner
				} else if x == 0 && y == r.ScreenHeight-1 {
					borderChar = '╚' // Bottom-left corner
				} else if x == r.ScreenWidth-1 && y == r.ScreenHeight-1 {
					borderChar = '╝' // Bottom-right corner
				} else if y == 0 || y == r.ScreenHeight-1 {
					borderChar = '═' // Horizontal border
				} else {
					borderChar = '║' // Vertical border
				}

				r.Pixels[y][x] = ColoredPixel{Char: borderChar, Color: "cyan"}
			} else if isGameTopLeft {
				r.Pixels[y][x] = ColoredPixel{Char: '╔', Color: "cyan"} // Top-left game field corner
			} else if isGameTopRight {
				r.Pixels[y][x] = ColoredPixel{Char: '╗', Color: "cyan"} // Top-right game field corner
			} else if isGameBottomLeft {
				r.Pixels[y][x] = ColoredPixel{Char: '╚', Color: "cyan"} // Bottom-left game field corner
			} else if isGameBottomRight {
				r.Pixels[y][x] = ColoredPixel{Char: '╝', Color: "cyan"} // Bottom-right game field corner
			} else if isLeftGameWall || isRightGameWall {
				r.Pixels[y][x] = ColoredPixel{Char: '║', Color: "cyan"} // Vertical game field walls
			} else if isGameFloor {
				r.Pixels[y][x] = ColoredPixel{Char: '═', Color: "cyan"} // Horizontal game field floor
			}
		}
	}
}

func (r *Renderer) Clear() {
	for y := 0; y < r.ScreenHeight; y++ {
		for x := 0; x < r.ScreenWidth; x++ {
			r.Pixels[y][x] = ColoredPixel{Char: ' ', Color: ""}
		}
	}
}

func (r *Renderer) RenderGame(game *Game) {
	// Clear the entire screen buffer
	r.Clear()

	// First build the border structure
	r.BuildBorder()

	game.UI.Draw(game)

	for _, block := range game.placedBlocks {
		r.RenderBlock(block, 0, 0)
	}

	if game.Player.CurrentPolymino != nil {
		r.DrawPolyomino(game.Player.CurrentPolymino)

		GetLoggerInstance().Log(fmt.Sprintf("Current Polyomino - Position: (%d, %d)",
			game.Player.CurrentPolymino.Position.X,
			game.Player.CurrentPolymino.Position.Y))
	}

	// For debugging - add the game field boundaries to logs
	gameRightWallScreenX := GameFieldStartX + (GameFieldWidth * BlockWidth)
	GetLoggerInstance().Log(fmt.Sprintf("Game boundaries: X=%d-%d, Y=%d-%d, Screen right wall at X=%d",
		GameFieldStartX, GameFieldEndX, GameFieldStartY, GameFieldEndY, gameRightWallScreenX))
	GetLoggerInstance().Log(fmt.Sprintf("Screen size: %dx%d", r.ScreenWidth, r.ScreenHeight))
}

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
)

func GetColorCode(color string) string {
	switch color {
	case "red":
		return ColorRed
	case "green":
		return ColorGreen
	case "yellow":
		return ColorYellow
	case "blue":
		return ColorBlue
	case "magenta":
		return ColorMagenta
	case "cyan":
		return ColorCyan
	case "white":
		return ColorWhite
	default:
		return ColorReset
	}
}

func (r *Renderer) RenderBlock(block Block, offsetX, offsetY int) {
	gameX := offsetX + block.Position.X
	gameY := offsetY + block.Position.Y

	screenX, screenY := r.GameToScreenCoordinates(gameX, gameY)

	if screenX >= 0 && screenX < r.ScreenWidth && screenY >= 0 && screenY < r.ScreenHeight {
		r.Pixels[screenY][screenX] = ColoredPixel{Char: '█', Color: block.Color}
		if screenX+1 < r.ScreenWidth {
			r.Pixels[screenY][screenX+1] = ColoredPixel{Char: '█', Color: block.Color}
		}
	}
}

func RenderMap(xCoord, yCoord int, pixelMap [][]rune) {
	for x := 0; x < len(pixelMap); x++ {
		for y := 0; y < len(pixelMap[0]); y++ {
			if xCoord+x < rendererInstance.ScreenWidth && yCoord+y < rendererInstance.ScreenHeight {
				if xCoord+x >= 0 && yCoord+y >= 0 {
					rendererInstance.Pixels[yCoord+y][xCoord+x] = ColoredPixel{Char: pixelMap[y][x], Color: ""}
				}
			}
		}
	}
}

func (r *Renderer) Render() {
	fmt.Print("\033c")

	for y := 0; y < r.ScreenHeight; y++ {
		for x := 0; x < r.ScreenWidth; x++ {
			pixel := r.Pixels[y][x]

			if pixel.Color != "" {
				colorCode := GetColorCode(pixel.Color)
				fmt.Print(colorCode + string(pixel.Char) + ColorReset)
			} else {
				fmt.Print(string(pixel.Char))
			}
		}
		fmt.Println()
	}

	GetLoggerInstance().PrintLogs(5)
}

func (r *Renderer) DrawPolyomino(polyomino *Polyomino) {
	for _, block := range polyomino.Blocks {
		gameX := polyomino.Position.X + block.Position.X
		gameY := polyomino.Position.Y + block.Position.Y

		screenX, screenY := r.GameToScreenCoordinates(gameX, gameY)

		if screenX >= 0 && screenX+1 < r.ScreenWidth && screenY >= 0 && screenY < r.ScreenHeight {
			r.Pixels[screenY][screenX] = ColoredPixel{Char: '█', Color: block.Color}
			r.Pixels[screenY][screenX+1] = ColoredPixel{Char: '█', Color: block.Color}
		}
	}
}

func (r *Renderer) GameToScreenCoordinates(gameX, gameY int) (int, int) {
	screenX := GameFieldStartX + (gameX * BlockWidth)
	screenY := GameFieldStartY + gameY
	return screenX, screenY
}

func (r *Renderer) ScreenToGameCoordinates(screenX, screenY int) (int, int) {
	gameX := (screenX - GameFieldStartX) / BlockWidth
	gameY := screenY - GameFieldStartY
	return gameX, gameY
}

func (r *Renderer) IsInGameArea(screenX, screenY int) bool {
	gameFieldEndScreenX := GameFieldStartX + (GameFieldWidth * BlockWidth)

	return screenX >= GameFieldStartX && screenX < gameFieldEndScreenX &&
		screenY >= GameFieldStartY && screenY < GameFieldEndY
}

func (r *Renderer) IsValidGameCoordinates(gameX, gameY int) bool {
	return gameX >= 0 && gameX < GameFieldWidth &&
		gameY >= 0 && gameY < GameFieldHeight
}
