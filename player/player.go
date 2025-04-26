package player

import (
	"github.com/iljarotar/synth/log"
	s "github.com/iljarotar/synth/synth"
)

type player struct {
	synth      *s.Synth
	sampleRate int
}

func NewPlayer(logger *log.Logger, filename string, sampleRate int) (*player, error) {
	p := &player{
		sampleRate: sampleRate,
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

func (p *player) UpdateSynth(synth *s.Synth) error {
	synth.Initialize(float64(p.sampleRate))
	p.synth = synth
	return nil
}
