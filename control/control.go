package control

import (
	"fmt"
	"math"

	"github.com/iljarotar/synth/audio"
	cfg "github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/log"
	"github.com/iljarotar/synth/synth"
	s "github.com/iljarotar/synth/synth"
)

type Control struct {
	logger                        *log.Logger
	config                        cfg.Config
	synth                         *s.Synth
	output                        chan audio.AudioOutput
	SynthDone                     chan bool
	autoStop                      chan bool
	maxOutput, lastNotifiedOutput float64
	overdriveWarningTriggeredAt   float64
	closing                       *bool
}

func NewControl(logger *log.Logger, config cfg.Config, output chan audio.AudioOutput, autoStop chan bool, closing *bool) (*Control, error) {
	var synth s.Synth
	err := synth.Initialize(config.SampleRate)
	if err != nil {
		return nil, err
	}

	ctl := &Control{
		logger:    logger,
		config:    config,
		synth:     &synth,
		output:    output,
		SynthDone: make(chan bool),
		autoStop:  autoStop,
		closing:   closing,
	}

	return ctl, nil
}

func (c *Control) Start() {
	outputChan := make(chan synth.Output)
	go c.synth.Play(outputChan)
	go c.receiveOutput(outputChan)
}

func (c *Control) LoadSynth(synth s.Synth) error {
	c.maxOutput = 0
	c.lastNotifiedOutput = 0

	err := synth.Initialize(c.config.SampleRate)
	if err != nil {
		return err
	}

	synth.Time += c.synth.Time

	*c.synth = synth

	return nil
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

func (c *Control) IncreaseVolume() {
	c.synth.IncreaseVolume()
}

func (c *Control) DecreaseVolume() {
	c.synth.DecreaseVolume()
	c.maxOutput = 0
	c.lastNotifiedOutput = 0
}

func (c *Control) GetVolume() float64 {
	return c.synth.VolumeMemory
}

func (c *Control) receiveOutput(outputChan <-chan synth.Output) {
	defer close(c.output)

	for out := range outputChan {
		c.logger.SendTime(out.Time)
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

	if c.maxOutput > 1 && c.maxOutput > c.lastNotifiedOutput && time-c.overdriveWarningTriggeredAt >= 0.5 {
		c.lastNotifiedOutput = c.maxOutput
		c.logger.ShowOverdriveWarning(true)
		c.logger.Warning(fmt.Sprintf("Output value %f", c.maxOutput))
		c.overdriveWarningTriggeredAt = time
	}
	if c.logger.OverdriveWarning && c.maxOutput <= 1 {
		c.logger.ShowOverdriveWarning(false)
	}
}

func (c *Control) checkDuration(time float64) {
	if c.config.Duration < 0 {
		return
	}
	duration := c.config.Duration - c.config.FadeOut
	if time < duration || *c.closing {
		return
	}
	c.autoStop <- true
}
