package control

import (
	"github.com/iljarotar/synth/audio"
	s "github.com/iljarotar/synth/synth"
	"gopkg.in/yaml.v2"
)

type Control struct {
	ctx   audio.Context
	Synth *s.Synth `yaml:"synth"`
}

func NewControl(ctx *audio.Context) *Control {
	return &Control{ctx: *ctx}
}

func (c *Control) Parse(data []byte) error {
	err := yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}

	c.Synth.Initialize()

	return nil
}

func (c *Control) Start() error {
	err := c.ctx.Start()
	if err != nil {
		return err
	}

	*c.Synth.Playing = true
	go c.Synth.Play(c.ctx.Input)

	return nil
}

func (c *Control) Stop() error {
	*c.Synth.Playing = false
	return c.ctx.Stop()
}
