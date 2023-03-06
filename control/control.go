package control

import (
	s "github.com/iljarotar/synth/synth"
)

type Control struct {
	synth *s.Synth
	input chan struct{ Left, Right float32 }
}

func NewControl(input chan struct{ Left, Right float32 }) *Control {
	var synth s.Synth
	synth.Initialize()
	ctl := &Control{synth: &synth, input: input}
	go ctl.synth.Play(ctl.input)
	return ctl
}

func (c *Control) LoadSynth(synth s.Synth) {
	synth.Initialize()
	synth.Phase = c.synth.Phase

	*c.synth = synth
}

func (c *Control) Close() {
	c.synth.Stop()
}

func (c *Control) Stop(fadeOut float64) {
	c.synth.FadeOut(fadeOut)
}

func (c *Control) Start(fadeIn float64) {
	// check if duration is set
	// if yes, start goroutine, that will send to the screens exit channel, once the duration is over
	c.synth.FadeIn(fadeIn)
}
