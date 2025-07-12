package game

import (
	"log"

	"github.com/eiannone/keyboard"
)

type EventHandler struct {
	Quit        chan bool
	InputEvents chan Event
}

type Event struct {
	Action    string
	Timestamp int64
}

var logger = GetLoggerInstance()

func NewEventHandler() *EventHandler {
	return &EventHandler{
		Quit:        make(chan bool),
		InputEvents: make(chan Event),
	}
}

func (e *EventHandler) Start() {
	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				logger.Log("Error reading key: " + err.Error())
				continue
			}

			switch key {
			case keyboard.KeyEsc:
				e.Quit <- true
				keyboard.Close()
				return
			case keyboard.KeyArrowDown:
				logger.Log("Key pressed: Down Arrow")
				e.InputEvents <- Event{Action: "down", Timestamp: 0}
			case keyboard.KeyArrowUp:
				logger.Log("Key pressed: Up Arrow")
				e.InputEvents <- Event{Action: "up"}
			case keyboard.KeyArrowLeft:
				logger.Log("Key pressed: Left Arrow")
				e.InputEvents <- Event{Action: "left"}
			case keyboard.KeyArrowRight:
				logger.Log("Key pressed: Right Arrow")
				e.InputEvents <- Event{Action: "right"}
			case keyboard.KeySpace:
				logger.Log("Key pressed: Space")
				e.InputEvents <- Event{Action: "space"}
			default:
				if char == 'c' || char == 'C' {
					logger.Log("Key pressed: C - Swap blocks")
					e.InputEvents <- Event{Action: "swap"}
				} else if char == 'd' || char == 'D' {
					logger.Log("Key pressed: D - Hard drop")
					e.InputEvents <- Event{Action: "hardDrop"}
				} else if char == 'r' || char == 'R' {
					logger.Log("Key pressed: R - Restart game")
					e.InputEvents <- Event{Action: "restart"}
				} else {
					log.Printf("Key pressed: %c", char)
				}
			}
		}
	}()
}

func (e *EventHandler) Stop() {
	close(e.Quit)
}

func (e *EventHandler) QuitChannel() <-chan bool {
	return e.Quit
}
