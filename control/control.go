package control

import (
	"fmt"
	"math"

	"github.com/iljarotar/synth/audio"
	cfg "github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/synth"
	s "github.com/iljarotar/synth/synth"
	"github.com/iljarotar/synth/ui"
)

type Control struct {
	logger                        *ui.Logger
	config                        cfg.Config
	synth                         *s.Synth
	output                        chan audio.AudioOutput
	SynthDone                     chan bool
	autoStop                      chan bool
	maxOutput, lastNotifiedOutput float64
	overdriveWarningTriggeredAt   float64
}

func NewControl(logger *ui.Logger, config cfg.Config, output chan audio.AudioOutput, autoStop chan bool) *Control {
	var synth s.Synth
	synth.Initialize(config.SampleRate)

	ctl := &Control{
		logger:    logger,
		config:    config,
		synth:     &synth,
		output:    output,
		SynthDone: make(chan bool),
		autoStop:  autoStop,
	}

	return ctl
}

func (c *Control) Start() {
	outputChan := make(chan synth.Output)
	go c.synth.Play(outputChan)
	go c.receiveOutput(outputChan)
}

func (c *Control) LoadSynth(synth s.Synth) {
	c.maxOutput = 0
	c.lastNotifiedOutput = 0
	synth.Initialize(c.config.SampleRate)
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

func (c *Control) receiveOutput(outputChan <-chan synth.Output) {
	defer close(c.output)

	for out := range outputChan {
		c.logTime(out.Time)
		c.checkDuration(out.Time)

		c.checkOverdrive(out.Mono, out.Time)

		c.output <- audio.AudioOutput{
			Left:  out.Left,
			Right: out.Right,
		}
	}
}

func (c *Control) checkOverdrive(output, time float64) {
	// only consider up to three decimals
	abs := math.Round(math.Abs(output)*1000) / 1000
	if abs > c.maxOutput {
		c.maxOutput = abs
	}

	if c.maxOutput >= 1 && c.maxOutput > c.lastNotifiedOutput && time-c.overdriveWarningTriggeredAt >= 0.5 {
		c.lastNotifiedOutput = c.maxOutput
		c.logger.ShowOverdriveWarning(true)
		c.logger.Warning(fmt.Sprintf("Output value %f", c.maxOutput))
		c.overdriveWarningTriggeredAt = time
	}
}

func (c *Control) checkDuration(time float64) {
	if c.config.Duration < 0 {
		return
	}
	duration := c.config.Duration - c.config.FadeOut
	if time < duration || ui.State.Closed {
		return
	}
	c.autoStop <- true
}

func (c *Control) logTime(time float64) {
	if isNextSecond(time) {
		c.logger.SendTime(int(time))
	}
}

func isNextSecond(time float64) bool {
	sec, _ := math.Modf(time)
	return sec > float64(ui.State.CurrentTime)
}
