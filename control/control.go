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

	return nil
}

func (c *Control) Close() error {
	return c.ctx.Close()
}

func (c *Control) Play() error {
	err := c.Start()
	if err != nil {
		return err
	}

	go c.Synth.Play(c.ctx.Input)
	return nil
}

func (c *Control) Stop() error {
	return c.ctx.Stop()
}
