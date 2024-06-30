package ui

import "fmt"

type logger struct {
	log                     chan string
	OverdriveWarningShowing bool
	overdriveWarning        chan bool
}

func (l *logger) Info(log string) {
	l.log <- fmt.Sprintf("\033[1;32m[INFO]\033[0m %s", log)
}

func (l *logger) Error(log string) {
	l.log <- fmt.Sprintf("\033[1;31m[ERROR]\033[0m %s", log)
}

func (l *logger) ShowOverdriveWarning() {
	l.OverdriveWarningShowing = true
	l.overdriveWarning <- true
}

var Logger = &logger{
	log:              make(chan string),
	overdriveWarning: make(chan bool),
}
