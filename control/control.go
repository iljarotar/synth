package control

import (
	"fmt"
	"math"

	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/log"
	"github.com/iljarotar/synth/synth"
)

type control struct {
	logger    *log.Logger
	config    *config.Config
	synth     *synth.Synth
	maxOutput float64
	notifyEnd chan<- bool
}

func NewControl(logger *log.Logger, c *config.Config) (*control, error) {
	p := &control{
		logger: logger,
		config: c,
	}

	return p, nil
}

func (c *control) ReadSample() [2]float64 {
	sample := [2]float64{}

	o := c.synth.GetOutput()
	sample[0] = o.Left
	sample[1] = o.Right

	c.logger.SendTime(o.Time)
	c.checkOutputLevel(o.Mono)

	if c.config.Duration > 0 && o.Time >= c.config.Duration && c.notifyEnd != nil {
		c.notifyEnd <- true
		c.notifyEnd = nil
	}

	return sample
}

func (c *control) WatchDuration(notifyEnd chan<- bool) {
	c.notifyEnd = notifyEnd
}

func (c *control) Stop(done chan<- bool, interrupt bool) {
	if c.synth == nil {
		done <- true
		return
	}

	fadeoutDone := make(chan bool)
	c.synth.NotifyFadeout(fadeoutDone)

	fadeout := c.config.FadeOut
	if interrupt {
		fadeout = 0.1
	}

	c.synth.FadeOut(fadeout)
	<-fadeoutDone
	done <- true
}

func (c *control) LoadSynth(synth *synth.Synth) error {
	err := synth.Initialize(float64(c.config.SampleRate))
	if err != nil {
		return err
	}

	if c.synth != nil {
		c.maxOutput = 0
		c.synth.Update(synth)
		return nil
	}

	c.synth = synth
	c.synth.FadeIn(c.config.FadeIn)
	return nil
}

func (c *control) IncreaseVolume() {
	c.synth.SetVolume(c.synth.Volume + 0.003)
}

func (c *control) DecreaseVolume() {
	c.synth.SetVolume(c.synth.Volume - 0.003)
	c.maxOutput = 0
}

func (c *control) Volume() float64 {
	return c.synth.VolumeMemory
}

func (c *control) checkOutputLevel(output float64) {
	// only consider up to three decimals
	abs := math.Round(math.Abs(output)*1000) / 1000
	if abs <= c.maxOutput {
		return
	}
	c.maxOutput = abs

	if c.maxOutput > 1.001 {
		c.logger.ShowVolumeWarning(true)
		c.logger.Warning(fmt.Sprintf("Output value %f", c.maxOutput))
	}
	if c.logger.State.VolumeWarning && c.maxOutput <= 1.001 {
		c.logger.ShowVolumeWarning(false)
	}
}
