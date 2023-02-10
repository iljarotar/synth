package control

import (
	s "github.com/iljarotar/synth/synth"
)

type Control struct {
	Synth *s.Synth
	input chan float32
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

	c.Synth.FadeOut(0.005)
	*c.Synth = synth
	c.Synth.FadeIn(0.005)
}

func (c *Control) Stop() {
	c.Synth.FadeOut(0.0001)
}

func (c *Control) Start() {
	go c.Synth.Play(c.input)
	c.Synth.FadeIn(0.005)
}
