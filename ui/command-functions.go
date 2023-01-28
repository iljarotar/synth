package ui

import (
	"os"
	"os/exec"

	"github.com/iljarotar/synth/config"
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

func setRootPath(done chan<- bool, args ...string) string {
	if len(args) != 1 {
		return "please specify exactly one root path"
	}

	c := config.Instance()
	*c.RootPath = args[0]

	return "root path set to " + *c.RootPath
}
