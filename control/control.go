package control

import (
	"time"

	"github.com/iljarotar/synth/config"
	s "github.com/iljarotar/synth/synth"
)

type Control struct {
	synth *s.Synth
	input chan struct{ Left, Right float32 }
	exit  chan bool
}

func NewControl(input chan struct{ Left, Right float32 }, exit chan bool) *Control {
	var synth s.Synth
	synth.Initialize()
	ctl := &Control{synth: &synth, input: input, exit: exit}
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
	if config.Config.Duration > 0 {
		go c.watchDuration()
	}
	c.synth.FadeIn(fadeIn)
}

func (c *Control) watchDuration() {
	time.Sleep(time.Duration(config.Config.Duration-config.Config.FadeOut) * time.Second)
	c.exit <- true
}
