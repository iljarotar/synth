package ui

import (
	"fmt"
)

const (
	labelInfo    = "[INFO]   "
	labelWarning = "[WARNING]"
	labelError   = "[ERROR]  "
)

type logger struct {
	log              chan string
	overdriveWarning chan bool
	time             chan string
}

func (l *logger) SendTime(time int) {
	State.CurrentTime = time
	l.time <- formatTime(time)
}

func (l *logger) Info(log string) {
	l.sendLog(log, labelInfo, colorGreenStrong)
}

func (l *logger) Warning(log string) {
	l.sendLog(log, labelWarning, colorOrangeStorng)
}

func (l *logger) Error(log string) {
	l.sendLog(log, labelError, colorRedStrong)
}

func (l *logger) ShowOverdriveWarning(limitExceeded bool) {
	State.ShowingOverdriveWarning = limitExceeded
	l.overdriveWarning <- limitExceeded
}

func (l *logger) sendLog(log, label string, labelColor color) {
	time := formatTime(State.CurrentTime)
	coloredLabel := fmt.Sprintf("%s", colored(label, labelColor))
	l.log <- fmt.Sprintf("[%s] %s %s", time, coloredLabel, log)
}

func colored(str string, col color) string {
	return fmt.Sprintf("%s%s%s", col, str, colorWhite)
}

func formatTime(time int) string {
	hours := time / 3600
	hoursString := fmt.Sprintf("%d", hours)
	if hours < 10 {
		hoursString = fmt.Sprintf("0%s", hoursString)
	}

	minutes := time / 60
	minutesString := fmt.Sprintf("%d", minutes)
	if minutes < 10 {
		minutesString = fmt.Sprintf("0%s", minutesString)
	}

	seconds := time % 60
	secondsString := fmt.Sprintf("%d", seconds)
	if seconds < 10 {
		secondsString = fmt.Sprintf("0%s", secondsString)
	}

	return fmt.Sprintf("%s:%s:%s", hoursString, minutesString, secondsString)
}

var Logger = &logger{
	log:              make(chan string),
	overdriveWarning: make(chan bool),
	time:             make(chan string),
}
