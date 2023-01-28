package ui

import (
	"fmt"
	"strings"

	"github.com/Songmu/prompter"
)

type UI struct {
	cmd cli
}

func NewUI() UI {
	u := UI{cmd: newCLI()}
	return u
}

func (u *UI) AcceptInput(done chan<- bool) {
	for {
		input := prompter.Prompt("", "")
		args := strings.Split(input, " ")
		resp := u.cmd.exec(args[0], done, args[1:]...)
		fmt.Println(resp)
	}
}

func (u *UI) ClearScreen() {
	cmd, ok := u.cmd.commands["clear"]
	if !ok {
		fmt.Println("something went wrong")
	}
	cmd(make(chan bool))
}
