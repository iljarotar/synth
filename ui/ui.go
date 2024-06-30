package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type UI struct {
	logger *Logger
	quit   chan bool
	input  chan string
	logs   []string
}

func NewUI(logger *Logger, quit chan bool) *UI {
	return &UI{logger: logger, quit: quit, input: make(chan string)}
}

func Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (u *UI) Enter(exit chan bool) {
	go u.read()
	u.resetScreen()

	for {
		select {
		case input := <-u.input:
			if input == "q" {
				u.quit <- true
				return
			} else {
				u.resetScreen()
			}
		case log := <-u.logger.log:
			u.logs = append(u.logs, log)
			u.resetScreen()
		case e := <-exit:
			if e == true {
				u.quit <- true
				return
			}
		}
	}
}

func (u *UI) read() {
	reader := bufio.NewReader(os.Stdin)

	for {
		in, _ := reader.ReadString('\n')
		u.input <- strings.TrimSpace(in)
	}
}

func (u *UI) resetScreen() {
	Clear()

	for i, log := range u.logs {
		fmt.Printf("[%d] %s\n", i+1, log)
	}
	if len(u.logs) > 0 {
		fmt.Print("\n")
	}
	fmt.Print("\033[1;34m Type 'q' to quit: \033[0m")
}
