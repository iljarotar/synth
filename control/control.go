package control

import (
	"github.com/iljarotar/synth/audio"
	s "github.com/iljarotar/synth/synth"
)

type Control struct {
	ctx         audio.Context
	Synth       *s.Synth
	Initialized bool
}

func NewControl(ctx *audio.Context) *Control {
	var synth s.Synth
	return &Control{ctx: *ctx, Synth: &synth, Initialized: false}
}

func (c *Control) LoadSynth(synth s.Synth) {
	synth.Initialize()
	*c.Synth = synth
	c.Initialized = true
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
