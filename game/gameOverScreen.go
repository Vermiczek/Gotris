package game

func (ui *Interface) DrawGameOverScreen() {
	gameOverX := GameFieldStartX + (GameFieldWidth * BlockWidth / 4)
	gameOverY := GameFieldStartY + (GameFieldHeight / 2)

	gameOverText := "GAME OVER"
	for i, char := range gameOverText {
		ui.renderer.Pixels[gameOverY][gameOverX+i] = ColoredPixel{Char: char, Color: "red"}
	}

	instructionsText := "Press R to restart"
	instructionsX := gameOverX - 2
	instructionsY := gameOverY + 2
	for i, char := range instructionsText {
		ui.renderer.Pixels[instructionsY][instructionsX+i] = ColoredPixel{Char: char, Color: "white"}
	}

	quitText := "ESC to quit"
	quitX := gameOverX + 2
	quitY := instructionsY + 1
	for i, char := range quitText {
		ui.renderer.Pixels[quitY][quitX+i] = ColoredPixel{Char: char, Color: "white"}
	}
}
