// FIX: removed output -> this component needs to be refactored completely

package control

import (
	"fmt"
	"math"

	cfg "github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/log"
	"github.com/iljarotar/synth/synth"
	s "github.com/iljarotar/synth/synth"
)

type Control struct {
	logger                        *log.Logger
	config                        cfg.Config
	synth                         *s.Synth
	SynthDone                     chan bool
	autoStop                      chan bool
	maxOutput, lastNotifiedOutput float64
	overdriveWarningTriggeredAt   float64
	closing                       *bool
}

func NewControl(logger *log.Logger, config cfg.Config, autoStop chan bool, closing *bool) (*Control, error) {
	var synth s.Synth
	err := synth.Initialize(float64(config.SampleRate))
	if err != nil {
		return nil, err
	}

	ctl := &Control{
		logger:    logger,
		config:    config,
		synth:     &synth,
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

	err := synth.Initialize(float64(c.config.SampleRate))
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
	for out := range outputChan {
		c.logger.SendTime(out.Time)
		c.checkDuration(out.Time)
		c.checkOverdrive(out.Mono, out.Time)
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
		c.logger.ShowVolumeWarning(true)
		c.logger.Warning(fmt.Sprintf("Output value %f", c.maxOutput))
		c.overdriveWarningTriggeredAt = time
	}
	if c.logger.VolumeWarning && c.maxOutput <= 1 {
		c.logger.ShowVolumeWarning(false)
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
