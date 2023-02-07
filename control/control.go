package control

import (
	s "github.com/iljarotar/synth/synth"
)

type Control struct {
	Synth   *s.Synth
	input   chan float32
	playing bool
}

func NewControl(input chan float32) *Control {
	var synth s.Synth
	synth.Initialize()
	ctl := &Control{Synth: &synth, input: input}
	return ctl
}

func (c *Control) LoadSynth(synth s.Synth) {
	synth.Initialize()
	synth.Phase = c.Synth.Phase

	c.Synth.FadeOut()
	*c.Synth = synth

	if c.playing {
		c.Synth.FadeIn()
	}
}

func (c *Control) Stop() {
	c.playing = false
	c.Synth.FadeOut()
}

func (c *Control) Start() {
	go c.Synth.Play(c.input)
	c.playing = true
	c.Synth.FadeIn()
}
