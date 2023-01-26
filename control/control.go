package control

import (
	"github.com/iljarotar/synth/audio"
	"github.com/iljarotar/synth/parser"
	s "github.com/iljarotar/synth/synth"
)

type Control struct {
	ctx     audio.Context
	Synth   *s.Synth
	playing *bool
}

func NewControl(ctx *audio.Context) *Control {
	return &Control{ctx: *ctx, playing: new(bool)}
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

	*c.playing = true
	go c.Synth.Play(c.ctx.Input, c.playing) // pass buffer instead

	return nil
}

func (c *Control) Stop() error {
	*c.playing = false
	return c.ctx.Stop()
}
