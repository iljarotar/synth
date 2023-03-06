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
	defer os.Stdin.Close()
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

	for i := range u.logs {
		fmt.Println(u.logs[i])
	}

	fmt.Print("type 'q' to quit: ")
}
