package control

import (
	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/log"
	"github.com/iljarotar/synth/synth"
)

type control struct {
	logger    *log.Logger
	config    *config.Config
	synth     *synth.Synth
	maxOutput float64
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

	return sample
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
		return c.synth.Update(synth)
	}

	c.synth = synth
	c.synth.FadeIn(c.config.FadeIn)
	return nil
}
