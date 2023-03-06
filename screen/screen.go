package screen

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Screen struct {
	logger *Logger
	quit   chan bool
	input  chan string
	logs   []string
}

func NewScreen(logger *Logger, quit chan bool) *Screen {
	return &Screen{logger: logger, quit: quit, input: make(chan string)}
}

func Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (s *Screen) Enter(exit chan bool) {
	defer os.Stdin.Close()
	go s.read()
	s.resetScreen()

	for {
		select {
		case input := <-s.input:
			if input == "q" {
				s.quit <- true
				return
			} else {
				s.resetScreen()
			}
		case log := <-s.logger.log:
			s.logs = append(s.logs, log)
			s.resetScreen()
		case e := <-exit:
			if e == true {
				s.quit <- true
				return
			}
		}
	}
}

func (s *Screen) read() {
	reader := bufio.NewReader(os.Stdin)

	for {
		in, _ := reader.ReadString('\n')
		s.input <- strings.TrimSpace(in)
	}
}

func (s *Screen) resetScreen() {
	Clear()

	for i := range s.logs {
		fmt.Println(s.logs[i])
	}

	fmt.Print("type 'q' to quit: ")
}
