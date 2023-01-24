package control

import (
	"github.com/iljarotar/synth/audio"
	s "github.com/iljarotar/synth/synth"
)

type Control struct {
	ctx   audio.Context
	synth *s.Synth
}

func NewControl(ctx audio.Context, synth *s.Synth) *Control {
	return &Control{ctx: ctx, synth: synth}
}

func (c *Control) Start() error {
	err := c.ctx.Start()
	if err != nil {
		return err
	}

	*c.synth.Playing = true
	go c.synth.Play(c.ctx.Input)

	return nil
}

func (c *Control) Stop() error {
	*c.synth.Playing = false
	return c.ctx.Stop()
}
