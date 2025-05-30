package log

import (
	"fmt"
	"math"
)

const (
	labelInfo    = "[INFO]   "
	labelWarning = "[WARNING]"
	labelError   = "[ERROR]  "
)

type State struct {
	Time          string
	VolumeWarning bool
}

type Logger struct {
	maxLogs          uint
	logs             []string
	State            State
	currentTime      int
	logSubscribers   []chan<- string
	stateSubscribers []chan<- State
}

func NewLogger(maxLogs uint) *Logger {
	return &Logger{maxLogs: maxLogs}
}

func (l *Logger) SubscribeToLogs(subscriber chan<- string) {
	l.logSubscribers = append(l.logSubscribers, subscriber)
}

func (l *Logger) SubscribeToState(subscriber chan<- State) {
	l.stateSubscribers = append(l.stateSubscribers, subscriber)
}

func (l *Logger) Info(log string) {
	l.sendLog(log, labelInfo, ColorGreenStrong)
}

func (l *Logger) Warning(log string) {
	l.sendLog(log, labelWarning, ColorOrangeStorng)
}

func (l *Logger) Error(log string) {
	l.sendLog(log, labelError, ColorRedStrong)
}

func (l *Logger) SendTime(time float64) {
	if l.isNextSecond(time) {
		seconds := int(time)
		l.currentTime = seconds
		l.State.Time = formatTime(seconds)
		l.sendState()
	}
}

func (l *Logger) ShowVolumeWarning(limitExceeded bool) {
	l.State.VolumeWarning = limitExceeded
	l.sendState()
}

func (l *Logger) sendState() {
	for _, s := range l.stateSubscribers {
		s <- l.State
	}
}

func (l *Logger) sendLog(log, label string, labelColor Color) {
	time := formatTime(l.currentTime)
	coloredLabel := fmt.Sprintf("%s", Colored(label, labelColor))

	for _, s := range l.logSubscribers {
		s <- fmt.Sprintf("[%s] %s %s", time, coloredLabel, log)
	}
}

func (l *Logger) isNextSecond(time float64) bool {
	sec, _ := math.Modf(time)
	return sec > float64(l.currentTime)
}

func Colored(str string, col Color) string {
	return fmt.Sprintf("%s%s%s", col, str, ColorWhite)
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
