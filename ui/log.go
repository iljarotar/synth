package ui

type Logger struct {
	log chan string
}

func NewLogger(log chan string) *Logger {
	return &Logger{log: log}
}

func (l *Logger) Log(log string) {
	l.log <- log
}
