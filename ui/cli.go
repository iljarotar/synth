package ui

import (
	c "github.com/iljarotar/synth/control"
	p "github.com/iljarotar/synth/parser"
)

type cli struct {
	commands map[string]cmdFunc
	config   cmdConfig
}

type cmdConfig struct {
	exit    chan<- bool
	control *c.Control
	parser  *p.Parser
}

func newCLI(config cmdConfig) cli {
	c := cli{config: config}
	c.commands = make(map[string]cmdFunc)
	c.addCommands()
	return c
}

func (c *cli) exec(input string, args ...string) string {
	cmd, ok := c.commands[input]
	if !ok {
		return "command not found"
	}

	return cmd(c.config, args...)
}

func (c *cli) addCommands() {
	c.commands["clear"] = clearCmd
	c.commands["c"] = clearCmd

	c.commands["exit"] = exitCmd
	c.commands["e"] = exitCmd

	c.commands["root"] = setRootPathCmd
	c.commands["r"] = setRootPathCmd

	c.commands["play"] = playCmd
	c.commands["p"] = playCmd

	c.commands["stop"] = stopCmd
	c.commands["s"] = stopCmd

	c.commands["load"] = loadCmd
	c.commands["l"] = loadCmd

	c.commands["apply"] = applyCmd
	c.commands["a"] = applyCmd
}
