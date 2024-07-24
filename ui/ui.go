package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type UI struct {
	quit  chan bool
	input chan string
	file  string
	logs  []string
	time  string
}

func NewUI(file string, quit chan bool) *UI {
	return &UI{
		quit:  quit,
		input: make(chan string),
		file:  file,
		time:  "00:00:00",
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
				ui.quit <- true
				return
			} else {
				ui.resetScreen()
			}
		case time := <-Logger.time:
			ui.time = time
			ui.updateTime()
		case log := <-Logger.log:
			ui.appendLog(log)
			ui.resetScreen()
		case <-Logger.overdriveWarning:
			ui.resetScreen()
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
	fmt.Printf("%s\n", ui.time)
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
