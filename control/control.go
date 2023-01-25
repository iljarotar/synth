package control

import (
	"github.com/iljarotar/synth/audio"
	"github.com/iljarotar/synth/parser"
	s "github.com/iljarotar/synth/synth"
)

type Control struct {
	ctx   audio.Context
	Synth *s.Synth
}

func NewControl(ctx *audio.Context) *Control {
	return &Control{ctx: *ctx}
}

func (c *Control) LoadSynth() error {
	var synth s.Synth
	err := parser.Parse(&synth)
	if err != nil {
		return err
	}

	synth.Initialize()
	c.Synth = &synth

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
