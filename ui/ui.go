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
}

func NewUI(file string, quit chan bool) *UI {
	return &UI{
		quit:  quit,
		input: make(chan string),
		file:  file,
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
				ui.quit <- true
				return
			} else {
				ui.resetScreen()
			}
		case log := <-Logger.log:
			ui.logs = append(ui.logs, log)
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
	fmt.Printf("%s%s", colored("Synth playing", COLOR_BLUE_STRONG), ui.file)
	LineBreaks(2)

	for i, log := range ui.logs {
		fmt.Printf(" [%d] %s", i+1, log)
		LineBreaks(1)
	}
	if len(ui.logs) > 0 {
		LineBreaks(1)
	}
	if Logger.ShowingOverdriveWarning {
		fmt.Printf("%s", colored("[WARNING] Volume exceeded 100%%", COLOR_ORANGE_STRONG))
		LineBreaks(2)
	}
	fmt.Printf("%s", colored("Type 'q' to quit:", COLOR_BLUE_STRONG))
}
