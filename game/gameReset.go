package game

func (g *Game) Reset() {
	g.placedBlocks = []Block{}
	g.lastDropTime = 0
	g.IsGameOver = false

	g.Player = NewPlayer()

	g.Scoring = NewScoringSystem()

	g.timer.Reset()

	GetLoggerInstance().Log("Game reset")
}
