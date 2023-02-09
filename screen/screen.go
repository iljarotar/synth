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
	done   chan bool
	input  chan string
	logs   []string
}

func NewScreen(logger *Logger, done chan bool) *Screen {
	return &Screen{logger: logger, done: done, input: make(chan string)}
}

func Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (s *Screen) Enter() {
	cancel := make(chan bool)
	go s.acceptInput(cancel)
	s.resetScreen()

	for {
		select {
		case input := <-s.input:
			if input == "q" {
				cancel <- true
				s.done <- true
			} else {
				s.resetScreen()
				cancel <- false
			}
		case log := <-s.logger.log:
			s.logs = append(s.logs, log)
			s.resetScreen()
		}
	}
}

func (s *Screen) acceptInput(cancel chan bool) {
	for {
		s.input <- prompt("")

		c := <-cancel
		if c == true {
			break
		}
	}
}

func (s *Screen) resetScreen() {
	Clear()

	for i := range s.logs {
		fmt.Println(s.logs[i])
	}

	fmt.Print("type 'q' to quit: ")
}

func prompt(label string) string {
	var input string
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(label)
		input, _ = reader.ReadString('\n')
		if input != "" {
			break
		}
	}

	return strings.TrimSpace(input)
}
