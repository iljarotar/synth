package player

import (
	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/log"
	"github.com/iljarotar/synth/synth"
)

type player struct {
	synth  *synth.Synth
	config *config.Config
}

func NewPlayer(logger *log.Logger, filename string, c *config.Config) (*player, error) {
	p := &player{
		config: c,
	}

	return p, nil
}

func (p *player) ReadSample() [2]float64 {
	sample := [2]float64{}
	o := p.synth.Next()
	sample[0] = o.Left
	sample[1] = o.Right
	return sample
}

func (p *player) LoadSynth(synth *synth.Synth) error {
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

func (p *player) Stop(done chan<- bool, interrupt bool) {
	if p.synth == nil {
		done <- true
		return
	}

	fadeoutDone := make(chan bool)
	p.synth.NotifyFadeout(fadeoutDone)

	fadeout := p.config.FadeOut
	if interrupt {
		fadeout = 0.05
	}

	p.synth.FadeOut(fadeout)
	<-fadeoutDone
	done <- true
	close(done)
}

func (p *player) updateSynth(synth *synth.Synth) {
	fadeoutDone := make(chan bool)
	p.synth.NotifyFadeout(fadeoutDone)

	p.synth.FadeOut(0.01)
	<-fadeoutDone
	p.synth = synth
	p.synth.FadeIn(0.01)
}
