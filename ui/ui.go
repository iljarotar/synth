package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/iljarotar/synth/log"
)

type (
	Signal string

	UI struct {
		logger     *log.Logger
		file       string
		signalChan chan<- Signal

		logs []string
		time string
	}

	Config struct {
		Logger     *log.Logger
		File       string
		Duration   float64
		SignalChan chan<- Signal
	}
)

const (
	SignalQuit      Signal = "quit"
	SignalInterrupt Signal = "interrupt"
)

func NewUI(c Config) *UI {
	return &UI{
		logger:     c.Logger,
		file:       c.File,
		signalChan: c.SignalChan,
		time:       "00:00:00",
	}
}

func LineBreaks(number int) {
	for range number {
		fmt.Print("\r\n")
	}
}

func (ui *UI) Enter() {
	go ui.read()
	ui.resetScreen()

	logChan := make(chan string)
	ui.logger.SubscribeToLogs(logChan)

	timeChan := make(chan string)
	ui.logger.SubscribeToTime(timeChan)

	for {
		select {
		case log := <-logChan:
			ui.appendLog(log)
			ui.resetScreen()

		case time := <-timeChan:
			if time != ui.time {
				ui.time = time
				ui.updateTime()
			}
		}
	}
}

func (ui *UI) read() {
	reader := bufio.NewReader(os.Stdin)

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			ui.logger.Error(fmt.Sprintf("failed to read input %v", err))
		}
		ui.handleInput(r)
	}
}

func (ui *UI) handleInput(r rune) {
	if r == rune(3) {
		ui.signalChan <- SignalInterrupt
		return
	}

	switch string(r) {
	case "q":
		ui.resetScreen()
		ui.signalChan <- SignalQuit
	}
}

func (ui *UI) resetScreen() {
	ui.clear()
	fmt.Printf("%s %s", log.Colored("Synth playing", log.ColorBlueStrong), ui.file)
	LineBreaks(2)

	for _, log := range ui.logs {
		fmt.Print(log + "\r\n")
	}
	if len(ui.logs) > 0 {
		LineBreaks(1)
	}
	fmt.Printf("%s ", ui.time)
	fmt.Print("Press 'q' to quit")
}

func (ui *UI) updateTime() {
	// using ANSI escape sequences:
	// \0337 to save current cursor location
	// \r to move cursor to beginning of line
	// \0338 to restore original cursor location
	fmt.Printf("\0337\r%s\0338", ui.time)
}

func (ui *UI) appendLog(log string) {
	logs := ui.logs
	logs = append(logs, log)
	if len(logs) > 10 {
		logs = logs[1:]
	}
	ui.logs = logs
}

func (ui *UI) clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		ui.logger.Error(err.Error())
	}
}
