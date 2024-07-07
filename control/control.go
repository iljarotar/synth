package control

import (
	s "github.com/iljarotar/synth/synth"
)

type Control struct {
	synth  *s.Synth
	output chan struct{ Left, Right float32 }
}

func NewControl(output chan struct{ Left, Right float32 }) *Control {
	var synth s.Synth
	synth.Initialize()
	ctl := &Control{synth: &synth, output: output}
	go ctl.synth.Play(ctl.output)
	return ctl
}

func (c *Control) LoadSynth(synth s.Synth) {
	synth.Initialize()
	synth.Time += c.synth.Time

	*c.synth = synth
}

func (c *Control) Close() {
	c.synth.Stop()
}

func (c *Control) Stop(fadeOut float64) {
	c.synth.FadeOut(fadeOut)
}

func (c *Control) Start(fadeIn float64) {
	c.synth.FadeIn(fadeIn)
}
