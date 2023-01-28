package ui

import (
	"os"
	"os/exec"
)

func clearFunc(done chan<- bool, args ...string) string {
	if len(args) > 0 {
		return "clear command doesn't expect any arguments"
	}

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	return ""
}

func exitFunc(done chan<- bool, args ...string) string {
	if len(args) > 0 {
		return "exit command doesn't expect any arguments"
	}

	done <- true
	return ""
}

func helpFunc(done chan<- bool, args ...string) string {
	if len(args) > 0 {
		return "help command doesn't expect any arguments"
	}

	return "some helpful information"
}
