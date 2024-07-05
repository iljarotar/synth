package ui

import "fmt"

type logger struct {
	log                     chan string
	ShowingOverdriveWarning bool
	overdriveWarning        chan bool
}

func (l *logger) Info(log string) {
	l.log <- fmt.Sprintf("\033[1;32m[INFO]\033[0m %s", log)
}

func (l *logger) Warning(log string) {
	l.log <- fmt.Sprintf("\033[1;33m[WARNING]\033[0m %s", log)
}

func (l *logger) Error(log string) {
	l.log <- fmt.Sprintf("\033[1;31m[ERROR]\033[0m %s", log)
}

func (l *logger) ShowOverdriveWarning(limitExceeded bool) {
	l.ShowingOverdriveWarning = limitExceeded
	l.overdriveWarning <- limitExceeded
}

var Logger = &logger{
	log:              make(chan string),
	overdriveWarning: make(chan bool),
}
