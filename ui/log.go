package ui

import (
	"fmt"
	"math"
)

const (
	labelInfo    = "[INFO]   "
	labelWarning = "[WARNING]"
	labelError   = "[ERROR]  "
)

type Logger struct {
	logChan              chan string
	overdriveWarningChan chan bool
	timeChan             chan string
	currentTime          int
}

func NewLogger() *Logger {
	return &Logger{
		logChan:              make(chan string),
		overdriveWarningChan: make(chan bool),
		timeChan:             make(chan string),
	}
}

func (l *Logger) SendTime(time float64) {
	if l.isNextSecond(time) {
		seconds := int(time)
		l.currentTime = seconds
		l.timeChan <- formatTime(seconds)
	}
}

func (l *Logger) Info(log string) {
	l.sendLog(log, labelInfo, colorGreenStrong)
}

func (l *Logger) Warning(log string) {
	l.sendLog(log, labelWarning, colorOrangeStorng)
}

func (l *Logger) Error(log string) {
	l.sendLog(log, labelError, colorRedStrong)
}

func (l *Logger) ShowOverdriveWarning(limitExceeded bool) {
	l.overdriveWarningChan <- limitExceeded
}

func (l *Logger) sendLog(log, label string, labelColor color) {
	time := formatTime(l.currentTime)
	coloredLabel := fmt.Sprintf("%s", colored(label, labelColor))
	l.logChan <- fmt.Sprintf("[%s] %s %s", time, coloredLabel, log)
}

func (l *Logger) isNextSecond(time float64) bool {
	sec, _ := math.Modf(time)
	return sec > float64(l.currentTime)
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

	minutes := time/60 - hours*60
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

// TODO:
// implement publish/subscribe mechanism
