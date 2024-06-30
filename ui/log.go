package ui

import "fmt"

type logger struct {
	log chan string
}

func (l *logger) Info(log string) {
	l.log <- fmt.Sprintf("\033[1;32m[INFO]\033[0m %s", log)
}

func (l *logger) Error(log string) {
	l.log <- fmt.Sprintf("\033[1;31m[ERROR]\033[0m %s", log)
}

var Logger = &logger{
	log: make(chan string),
}
