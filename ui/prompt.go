package ui

import (
	"fmt"
	"strings"

	"github.com/Songmu/prompter"
	c "github.com/iljarotar/synth/control"
)

type UI struct {
	cli cli
}

func NewUI(ctl *c.Control, exit chan<- bool) UI {
	config := cmdConfig{control: ctl, exit: exit}
	return UI{cli: newCLI(config)}
}

func (ui *UI) AcceptInput() {
	for {
		input := prompter.Prompt("", "")
		args := strings.Split(input, " ")
		resp := ui.cli.exec(args[0], args[1:]...)
		fmt.Println(resp)
	}
}

func (ui *UI) ClearScreen() {
	cmd, ok := ui.cli.commands["clear"]
	if !ok {
		fmt.Println("something went wrong")
	}
	cmd(ui.cli.config)
}
