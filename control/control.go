package control

import (
	"github.com/iljarotar/synth/audio"
	s "github.com/iljarotar/synth/synth"
)

type Control struct {
	ctx         audio.Context
	Synth       *s.Synth
	Initialized bool
	playing     bool
}

func NewControl(ctx *audio.Context) *Control {
	var synth s.Synth
	synth.Initialize()
	ctl := &Control{ctx: *ctx, Synth: &synth, Initialized: false, playing: false}
	ctl.Start()
	return ctl
}

func (c *Control) LoadSynth(synth s.Synth) {
	synth.Initialize()
	c.Synth.FadeOut()
	*c.Synth = synth

	if c.playing {
		c.Synth.FadeIn()
	}

	c.Initialized = true
}

func (c *Control) Play() {
	c.Synth.FadeIn()
	c.playing = true
}

func (c *Control) Stop() {
	c.Synth.FadeOut()
	c.playing = false
}

func (c *Control) Start() error {
	err := c.ctx.Start()
	if err != nil {
		return err
	}

	go c.Synth.Play(c.ctx.Input)

	return nil
}

func (c *Control) Close() error {
	defer c.ctx.Close()

	err := c.ctx.Stop()
	if err != nil {
		return err
	}

	return nil
}
