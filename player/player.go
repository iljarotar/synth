package player

import (
	"github.com/iljarotar/synth/module"
	s "github.com/iljarotar/synth/synth"
)

type player struct {
	synth *s.Synth
}

func NewPlayer(filename string, sampleRate int) (*player, error) {
	// TODO: read file and initialize synth but for now test with some hard-coded config
	s := &s.Synth{
		Volume: 1,
		Out:    []string{"o1"},
		Oscillators: []*module.Oscillator{
			{
				Name: "o1",
				Type: "Sine",
				Freq: module.Input{
					Val: 220,
				},
				Amp: module.Input{
					Val: 1,
				},
				Pan: module.Input{},
			},
		},
	}
	err := s.Initialize(44100)
	if err != nil {
		return nil, err
	}

	p := &player{
		synth: s,
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
