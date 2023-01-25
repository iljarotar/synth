package ui

import (
	"os"
	"os/exec"

	"github.com/Songmu/prompter"
)

func AcceptInput(done chan<- bool) {
	for input := prompter.Prompt(">", ""); input != "exit"; {
		input = prompter.Prompt(">", "")
	}
	done <- true
}

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
