package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type UI struct {
	logger   *Logger
	quit     chan bool
	input    chan string
	autoStop chan bool
	file     string
	logs     []string
	time     string
	duration float64
}

func NewUI(logger *Logger, file string, quit chan bool, autoStop chan bool, duration float64) *UI {
	return &UI{
		logger:   logger,
		quit:     quit,
		autoStop: autoStop,
		input:    make(chan string),
		file:     file,
		time:     "00:00:00",
		duration: duration,
	}
}

func Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func LineBreaks(number int) {
	for i := 0; i < number; i++ {
		fmt.Print("\n")
	}
}

func (ui *UI) Enter() {
	go ui.read()
	ui.resetScreen()

	for {
		select {
		case input := <-ui.input:
			if input == "q" {
				State.Closed = true
				ui.resetScreen()
				ui.quit <- true
			} else {
				ui.resetScreen()
			}
		case time := <-ui.logger.time:
			ui.time = time
			ui.updateTime()
		case log := <-ui.logger.log:
			ui.appendLog(log)
			ui.resetScreen()
		case <-ui.logger.overdriveWarning:
			ui.resetScreen()
		case <-ui.autoStop:
			State.Closed = true
			ui.quit <- true
		}
	}
}

func (ui *UI) read() {
	reader := bufio.NewReader(os.Stdin)

	for {
		in, _ := reader.ReadString('\n')
		ui.input <- strings.TrimSpace(in)
	}
}

func (ui *UI) resetScreen() {
	Clear()
	LineBreaks(1)
	fmt.Printf("%s %s", colored("Synth playing", colorBlueStrong), ui.file)
	LineBreaks(2)

	for _, log := range ui.logs {
		fmt.Println(log)
	}
	if len(ui.logs) > 0 {
		LineBreaks(1)
	}
	if State.ShowingOverdriveWarning {
		fmt.Printf("%s", colored("[WARNING] Volume exceeded 100%%", colorOrangeStorng))
		LineBreaks(2)
	}
	fmt.Printf("%s", ui.time)
	if ui.duration >= 0 {
		fmt.Printf(" - automatically stopping after %fs", ui.duration)
	}
	LineBreaks(1)
	fmt.Printf("%s ", colored("Type 'q' to quit:", colorBlueStrong))
}

func (ui *UI) updateTime() {
	// using ANSI escape sequences:
	// \0337 to save current cursor location
	// \033[1A to move cursor up one line
	// \r to move cursor to beginning of line
	// \0338 to restore original cursor location
	fmt.Printf("\0337\033[1A\r%s\0338", ui.time)
}

func (ui *UI) appendLog(log string) {
	logs := ui.logs
	logs = append(logs, log)
	if len(logs) > 10 {
		logs = logs[1:]
	}
	ui.logs = logs
}
