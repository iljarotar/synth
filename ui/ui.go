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
	logs  []string
}

func NewUI(quit chan bool) *UI {
	return &UI{quit: quit, input: make(chan string)}
}

func Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
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

	for i, log := range ui.logs {
		fmt.Printf("[%d] %s\n", i+1, log)
	}
	if len(ui.logs) > 0 {
		fmt.Print("\n")
	}
	fmt.Print("\033[1;34m Type 'q' to quit: \033[0m")
}
