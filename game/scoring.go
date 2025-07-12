package game

const (
	PointsPerLine       = 100
	PointsPerLineTetris = 400
)

type ScoringSystem struct {
	Score        int
	LinesCleared int
	Level        int
}

func NewScoringSystem() *ScoringSystem {
	return &ScoringSystem{
		Score:        0,
		LinesCleared: 0,
		Level:        1,
	}
}

func (s *ScoringSystem) AddScore(points int) {
	s.Score += points * s.Level
}

func (s *ScoringSystem) AddLines(lines int) {
	if lines <= 0 {
		return
	}

	s.LinesCleared += lines

	s.Level = (s.LinesCleared / 10) + 1

	if lines == 4 {
		s.AddScore(PointsPerLineTetris)
	} else {
		s.AddScore(PointsPerLine * lines)
	}
}

func (s *ScoringSystem) GetDropSpeed() int64 {
	speed := 1000 - ((s.Level - 1) * 100)
	if speed < 100 {
		speed = 100
	}
	return int64(speed)
}
