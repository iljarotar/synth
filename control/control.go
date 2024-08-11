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
	logger      *ui.Logger
	config      cfg.Config
	synth       *s.Synth
	output      chan audio.AudioOutput
	SynthDone   chan bool
	autoStop    chan bool
	currentTime float64
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
	for out := range outputChan {
		c.currentTime = out.Time
		c.logTime(out.Time)
		c.checkDuration()

		c.checkOverdrive(out.Mono)

		c.output <- audio.AudioOutput{
			Left:  out.Left,
			Right: out.Right,
		}
	}
}

func (c *Control) checkOverdrive(output float64) {
	if math.Abs(output) >= 1.00001 && !ui.State.ShowingOverdriveWarning {
		c.logger.ShowOverdriveWarning(true)
		c.logger.Warning(fmt.Sprintf("Output value %f", output))
	}
}

func (c *Control) checkDuration() {
	if c.config.Duration < 0 {
		return
	}
	duration := c.config.Duration - c.config.FadeOut
	if c.currentTime < duration || ui.State.Closed {
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
