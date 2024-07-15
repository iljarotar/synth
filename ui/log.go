package ui

import (
	"fmt"
)

type logger struct {
	log                     chan string
	ShowingOverdriveWarning bool
	overdriveWarning        chan bool
}

func (l *logger) Info(log string) {
	if State.Closed {
		return
	}
	l.log <- fmt.Sprintf("%s %s", colored("[INFO]", COLOR_GREEN_STRONG), log)
}

func (l *logger) Warning(log string) {
	if State.Closed {
		return
	}
	l.log <- fmt.Sprintf("%s %s", colored("[WARNING]", COLOR_ORANGE_STRONG), log)
}

func (l *logger) Error(log string) {
	if State.Closed {
		return
	}
	l.log <- fmt.Sprintf("%s %s", colored("[EROOR]", COLOR_RED_STRONG), log)
}

func (l *logger) ShowOverdriveWarning(limitExceeded bool) {
	if State.Closed {
		return
	}
	l.ShowingOverdriveWarning = limitExceeded
	l.overdriveWarning <- limitExceeded
}

func colored(str string, col color) string {
	return fmt.Sprintf("%s %s %s", col, str, COLOR_WHITE)
}

var Logger = &logger{
	log:              make(chan string),
	overdriveWarning: make(chan bool),
}
