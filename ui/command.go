package ui

type cmdFunc func(done chan<- bool, args ...string) string

type cli struct {
	commands map[string]cmdFunc
}

func newCLI() cli {
	c := cli{}
	c.commands = make(map[string]cmdFunc)
	c.commands["clear"] = clearFunc
	c.commands["c"] = clearFunc

	c.commands["exit"] = exitFunc
	c.commands["e"] = exitFunc

	c.commands["help"] = helpFunc
	c.commands["h"] = helpFunc

	c.commands["root"] = setRootPath
	c.commands["r"] = setRootPath
	return c
}

func (c *cli) exec(input string, done chan<- bool, args ...string) string {
	cmd, ok := c.commands[input]
	if !ok {
		return "command not found"
	}

	return cmd(done, args...)
}
