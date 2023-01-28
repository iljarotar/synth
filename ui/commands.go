package ui

import (
	"os"
	"os/exec"

	c "github.com/iljarotar/synth/config"
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

	*c.Instance().RootPath = args[0]
	return "root path set to " + *c.Instance().RootPath
}

func playCmd(config cmdConfig, args ...string) string {
	if len(args) > 0 {
		return "play command doesn't expect any arguments"
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
