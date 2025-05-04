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
}

func NewControl(logger *log.Logger, c *config.Config) (*control, error) {
	p := &control{
		logger: logger,
		config: c,
	}

	return p, nil
}

func (p *control) ReadSample() [2]float64 {
	sample := [2]float64{}

	o := p.synth.Next()
	sample[0] = o.Left
	sample[1] = o.Right

	p.logger.SendTime(o.Time)
	p.checkOutputLevel(o.Mono)

	return sample
}

func (p *control) Stop(done chan<- bool, interrupt bool) {
	if p.synth == nil {
		done <- true
		return
	}

	fadeoutDone := make(chan bool)
	p.synth.NotifyFadeout(fadeoutDone)

	fadeout := p.config.FadeOut
	if interrupt {
		fadeout = 0.1
	}

	p.synth.FadeOut(fadeout)
	<-fadeoutDone
	done <- true
	close(done)
}

func (p *control) LoadSynth(synth *synth.Synth) error {
	err := synth.Initialize(float64(p.config.SampleRate))
	if err != nil {
		return err
	}

	if p.synth != nil {
		p.updateSynth(synth)
		return nil
	}

	p.synth = synth
	p.synth.FadeIn(p.config.FadeIn)
	return nil
}

func (p *control) IncreaseVolume() {
	p.synth.SetVolume(p.synth.Volume + 0.003)
}

func (p *control) DecreaseVolume() {
	p.synth.SetVolume(p.synth.Volume - 0.003)
}

func (p *control) Volume() float64 {
	return p.synth.VolumeMemory
}

func (p *control) updateSynth(synth *synth.Synth) {
	fadeoutDone := make(chan bool)
	p.synth.NotifyFadeout(fadeoutDone)
	p.synth.FadeOut(0.01)
	<-fadeoutDone

	p.maxOutput = 0
	synth.Time = p.synth.Time
	p.synth = synth
	p.synth.FadeIn(0.01)
}

func (p *control) checkOutputLevel(output float64) {
	// only consider up to three decimals
	abs := math.Round(math.Abs(output)*1000) / 1000
	if abs <= p.maxOutput {
		return
	}
	p.maxOutput = abs

	if p.maxOutput > 1.001 {
		p.logger.ShowVolumeWarning(true)
		p.logger.Warning(fmt.Sprintf("Output value %f", p.maxOutput))
	}
	if p.logger.State.VolumeWarning && p.maxOutput <= 1.001 {
		p.logger.ShowVolumeWarning(false)
	}
}
