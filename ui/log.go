package ui

type logger struct {
	log chan string
}

func (l *logger) Log(log string) {
	l.log <- log
}

var Logger = &logger{
	log: make(chan string),
}
