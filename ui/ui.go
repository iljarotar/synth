package ui

import (
	"fmt"
	"strings"

	"github.com/Songmu/prompter"
	c "github.com/iljarotar/synth/control"
	p "github.com/iljarotar/synth/parser"
)

type UI struct {
	cli cli
}

func NewUI(ctl *c.Control, exit chan<- bool) UI {
	parser := p.NewParser()
	config := cmdConfig{control: ctl, exit: exit, parser: parser}
	return UI{cli: newCLI(config)}
}

func (ui *UI) AcceptInput() {
	for {
		input := prompter.Prompt("", "")
		args := strings.Split(input, " ")
		resp := ui.cli.exec(args[0], args[1:]...)

		if resp != "" {
			fmt.Println(resp)
		}
	}
}

func (ui *UI) ClearScreen(msg ...string) {
	cmd, ok := ui.cli.commands["clear"]
	if !ok {
		fmt.Println("something went wrong")
	}
	cmd(ui.cli.config)

	if len(msg) > 0 && msg[0] != "" {
		fmt.Println(msg[0])
	}
}

func (ui *UI) PrintMenu() {
	fmt.Println(menu)
}
