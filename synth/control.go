package synth

import (
	"math"

	"github.com/iljarotar/synth/audio"
	cfg "github.com/iljarotar/synth/config"
)

type Control struct {
	config                        cfg.Config
	synth                         *Synth
	output                        chan audio.AudioOutput
	synthDone                     chan bool
	autoStop                      chan bool
	maxOutput, lastNotifiedOutput float64
	overdriveWarningTriggeredAt   float64
	closing                       *bool
}

func NewControl(synth *Synth, config cfg.Config, output chan audio.AudioOutput, autoStop chan bool) (*Control, error) {
	err := synth.initialize(config.SampleRate)
	if err != nil {
		return nil, err
	}

	ctl := &Control{
		config:    config,
		synth:     synth,
		output:    output,
		synthDone: make(chan bool),
		autoStop:  autoStop,
	}

	return ctl, nil
}

func (c *Control) GetVolume() float64 {
	return c.synth.Volume
}

func (c *Control) IncreaseVolume() {
	vol := c.synth.Volume + 0.02
	if vol > maxVolume {
		vol = maxVolume
	}
	c.synth.volumeMemory = vol
	c.synth.Volume = vol
}

func (c *Control) DecreaseVolume() {
	vol := c.synth.Volume - 0.02
	if vol < 0 {
		vol = 0
	}
	c.synth.volumeMemory = vol
	c.synth.Volume = vol
}

func (c *Control) Start() {
	outputChan := make(chan Output)
	c.synth.active = true
	go c.synth.play(outputChan)
	go c.receiveOutput(outputChan)
	c.synth.startFading(FadeDirectionIn, c.config.FadeIn)
}

func (c *Control) Stop() {
	c.synth.notifyFadeOutDone(c.synthDone)
	c.synth.startFading(FadeDirectionOut, c.config.FadeOut)
	<-c.synthDone
	c.synth.active = false
}

func (c *Control) checkOverdrive(output, time float64) {
	// only consider up to three decimals
	abs := math.Round(math.Abs(output)*1000) / 1000
	if abs > c.maxOutput {
		c.maxOutput = abs
	}

	if c.maxOutput > 1 && c.maxOutput > c.lastNotifiedOutput && time-c.overdriveWarningTriggeredAt >= 0.5 {
		c.lastNotifiedOutput = c.maxOutput
		// TODO: overdrive warning
		c.overdriveWarningTriggeredAt = time
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

func (c *Control) receiveOutput(outputChan <-chan Output) {
	defer close(c.output)

	for out := range outputChan {
		c.checkDuration(out.Time)
		c.checkOverdrive(out.Mono, out.Time)

		c.output <- audio.AudioOutput{
			Left:  out.Left,
			Right: out.Right,
		}
	}
}
