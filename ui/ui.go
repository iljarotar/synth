package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/iljarotar/synth/log"
)

type UI struct {
	logger   *log.Logger
	file     string
	duration float64
	quitChan chan<- bool

	input             chan string
	logs              []string
	time              string
	showVolumeWarning bool
}

type Config struct {
	Logger   *log.Logger
	File     string
	Duration float64
	QuitChan chan<- bool
}

func NewUI(c Config) *UI {
	return &UI{
		logger:   c.Logger,
		file:     c.File,
		duration: c.Duration,
		quitChan: c.QuitChan,
		input:    make(chan string),
		time:     "00:00:00",
	}
}

func Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func lineBreaks(number int) {
	for range number {
		fmt.Print("\r\n")
	}
}

func (ui *UI) Enter() {
	go ui.read()
	ui.resetScreen()

	logChan := make(chan string)
	ui.logger.SubscribeToLogs(logChan)

	stateChan := make(chan log.State)
	ui.logger.SubscribeToState(stateChan)

	for {
		select {
		case input := <-ui.input:
			switch input {
			case "q":
				ui.resetScreen()
				ui.quitChan <- true
			case "d":
				// TODO: increase volume
				ui.resetScreen()
			case "s":
				// TODO: decrease volume
				ui.resetScreen()
			}

		case log := <-logChan:
			ui.appendLog(log)
			ui.resetScreen()

		case state := <-stateChan:
			if state.VolumeWarning != ui.showVolumeWarning {
				ui.showVolumeWarning = state.VolumeWarning
				ui.resetScreen()
			}
			if state.Time != ui.time {
				ui.time = state.Time
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
		if r == rune(3) {
			// TODO: interrupt?
		}
		ui.input <- string(r)
	}
}

func (ui *UI) resetScreen() {
	Clear()
	fmt.Printf("%s %s", log.Colored("Synth playing", log.ColorBlueStrong), ui.file)
	lineBreaks(1)
	fmt.Printf("%s %s", log.Colored("Volume", log.ColorBlueStrong), fmt.Sprintf("%v", 1)) // TODO: get volume
	lineBreaks(2)

	for _, log := range ui.logs {
		fmt.Print(log + "\r\n")
	}
	if len(ui.logs) > 0 {
		lineBreaks(1)
	}
	if ui.showVolumeWarning {
		fmt.Printf("%s", log.Colored("[WARNING] Volume exceeded 100%%", log.ColorOrangeStorng))
		lineBreaks(2)
	}
	fmt.Printf("%s", ui.time)
	if ui.duration >= 0 {
		fmt.Printf(" - automatically stopping after %fs", ui.duration)
	}
	lineBreaks(2)
	fmt.Printf("%s ", log.Colored("Keybindings", log.ColorBlueStrong))
	lineBreaks(1)
	fmt.Print("q: quit")
	lineBreaks(1)
	fmt.Print("d: raise volume")
	lineBreaks(1)
	fmt.Print("s: reduce volume")
}

func (ui *UI) updateTime() {
	// using ANSI escape sequences:
	// \0337 to save current cursor location
	// \033[1A to move cursor up one line
	// \r to move cursor to beginning of line
	// \0338 to restore original cursor location
	fmt.Printf("\0337\033[5A\r%s\0338", ui.time)
}

func (ui *UI) appendLog(log string) {
	logs := ui.logs
	logs = append(logs, log)
	if len(logs) > 10 {
		logs = logs[1:]
	}
	ui.logs = logs
}
