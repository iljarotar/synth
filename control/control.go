package control

import (
	s "github.com/iljarotar/synth/synth"
)

type Control struct {
	synth     *s.Synth
	output    chan struct{ Left, Right float32 }
	SynthDone chan bool
}

func NewControl(output chan struct{ Left, Right float32 }) *Control {
	var synth s.Synth
	synth.Initialize()
	ctl := &Control{synth: &synth, output: output, SynthDone: make(chan bool)}
	go ctl.synth.Play(ctl.output)
	return ctl
}

func (c *Control) LoadSynth(synth s.Synth) {
	synth.Initialize()
	synth.Time += c.synth.Time

	*c.synth = synth
}

func (c *Control) Stop(fadeOut float64) {
	c.FadeOut(fadeOut, c.SynthDone)
}

func (c *Control) StopSynth() {
	c.synth.Stop()
}

func (c *Control) FadeIn(fadeIn float64) {
	c.synth.Fade(s.FadeDirectionIn, fadeIn)
}

func (c *Control) FadeOut(fadeOut float64, notifyDone chan bool) {
	c.synth.NotifyFadeOutDone(notifyDone)
	c.synth.Fade(s.FadeDirectionOut, fadeOut)
}
