package game

import "time"

type GameTimer struct {
	elapsed    int64
	startTime  time.Time
	lastTick   time.Time
	tickRate   time.Duration
	isPaused   bool
	pausedTime int64
}

func NewGameTimer() *GameTimer {
	return &GameTimer{
		elapsed:    0,
		startTime:  time.Now(),
		lastTick:   time.Now(),
		tickRate:   100 * time.Millisecond,
		isPaused:   false,
		pausedTime: 0,
	}
}

func (t *GameTimer) Update() {
	if t.isPaused {
		return
	}

	now := time.Now()
	elapsed := now.Sub(t.lastTick)
	t.elapsed += elapsed.Milliseconds()
	t.lastTick = now
}

func (t *GameTimer) Reset() {
	now := time.Now()
	t.elapsed = 0
	t.startTime = now
	t.lastTick = now
	t.pausedTime = 0
}

func (t *GameTimer) Pause() {
	if !t.isPaused {
		t.isPaused = true
		t.pausedTime = time.Now().UnixMilli()
	}
}

func (t *GameTimer) Resume() {
	if t.isPaused {
		t.isPaused = false
		now := time.Now()
		t.lastTick = now
	}
}
