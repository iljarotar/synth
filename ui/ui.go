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

func (ui *UI) Enter(exit chan bool) {
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
		case e := <-exit:
			if e == true {
				ui.quit <- true
				return
			}
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
	fmt.Printf("\033[1;34m Synth playing \033[0m %s", ui.file)
	LineBreaks(2)

	for i, log := range ui.logs {
		fmt.Printf(" [%d] %s", i+1, log)
		LineBreaks(1)
	}
	if len(ui.logs) > 0 {
		LineBreaks(1)
	}
	if Logger.ShowingOverdriveWarning {
		fmt.Printf("\033[1;33m [WARNING] Volume exceeded 100%% \033[0m")
		LineBreaks(2)
	}
	fmt.Print("\033[1;34m Type 'q' to quit: \033[0m")
}
