package control

import (
	"math"

	"github.com/iljarotar/synth/config"
	s "github.com/iljarotar/synth/synth"
	"github.com/iljarotar/synth/ui"
)

type Control struct {
	synth       *s.Synth
	output      chan struct{ Left, Right float32 }
	SynthDone   chan bool
	autoStop    chan bool
	reportTime  chan float64
	currentTime float64
}

func NewControl(output chan struct{ Left, Right float32 }, autoStop chan bool) *Control {
	var synth s.Synth
	synth.Initialize()
	reportTime := make(chan float64)

	ctl := &Control{
		synth:      &synth,
		output:     output,
		SynthDone:  make(chan bool),
		reportTime: reportTime,
		autoStop:   autoStop,
	}

	return ctl
}

func (c *Control) Start() {
	go c.synth.Play(c.output, c.reportTime)
	go c.observeTime()
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

func (c *Control) observeTime() {
	for time := range c.reportTime {
		c.currentTime = time
		logTime(time)
		c.checkDuration()
	}
}

func (c *Control) checkDuration() {
	if config.Config.Duration < 0 {
		return
	}
	duration := config.Config.Duration - config.Config.FadeOut
	if c.currentTime < duration || ui.State.Closed {
		return
	}
	c.autoStop <- true
}

func logTime(time float64) {
	if isNextSecond(time) {
		ui.Logger.SendTime(int(time))
	}
}

func isNextSecond(time float64) bool {
	sec, _ := math.Modf(time)
	return sec > float64(ui.State.CurrentTime)
}
