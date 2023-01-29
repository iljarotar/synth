package ui

import (
	"os"
	"os/exec"

	c "github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/synth"
)

type cmdFunc func(config cmdConfig, args ...string) string

func clearCmd(config cmdConfig, args ...string) string {
	if len(args) > 0 {
		return "clear command doesn't expect any arguments"
	}

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	return ""
}

func exitCmd(config cmdConfig, args ...string) string {
	if len(args) > 0 {
		return "exit command doesn't expect any arguments"
	}

	config.exit <- true
	return ""
}

func setRootPathCmd(config cmdConfig, args ...string) string {
	if len(args) != 1 {
		return "please specify exactly one root path"
	}

	err := config.parser.SetRootPath(args[0])
	if err != nil {
		return "could not set root path"
	}

	path := c.Instance().RootPath
	return "root path set to " + path
}

func playCmd(config cmdConfig, args ...string) string {
	if len(args) > 0 {
		return "play command doesn't expect any arguments"
	}

	if config.control.Initialized == false {
		return "don't know what to play. please load a file first"
	}

	config.control.Play()
	return ""
}

func stopCmd(config cmdConfig, args ...string) string {
	if len(args) > 0 {
		return "stop command doesn't expect any arguments"
	}

	config.control.Stop()
	return ""
}

func loadCmd(config cmdConfig, args ...string) string {
	if len(args) != 1 {
		return "please specify exactly one file to load"
	}

	var s synth.Synth
	err := config.parser.Load(args[0], &s)
	if err != nil {
		return err.Error()
	}

	config.control.LoadSynth(s)
	return ""
}

func applyCmd(config cmdConfig, args ...string) string {
	if len(args) > 0 {
		return "apply command doesn't expect any arguments"
	}

	var s synth.Synth
	err := config.parser.LoadLastOpened(&s)
	if err != nil {
		return err.Error()
	}

	config.control.LoadSynth(s)
	return ""
}

func helpCmd(config cmdConfig, args ...string) string {
	if len(args) > 0 {
		return "help command doesn't expect any arguments"
	}

	return menu
}
