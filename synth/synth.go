package synth

import (
	"time"

	"github.com/iljarotar/synth/config"
	w "github.com/iljarotar/synth/wavetable"
)

type Synth struct {
	Gain                    float64 `yaml:"gain"`
	gainMemory, Phase, step float64
	WaveTables              []*w.WaveTable `yaml:"wavetables"`
}

func (s *Synth) Initialize() {
	c := config.Instance()
	s.step = 1 / c.SampleRate

	if s.Gain == 0 {
		s.Gain = 1
	}
	s.gainMemory = s.Gain
	s.Gain = 0 // start muted

	for i := range s.WaveTables {
		s.WaveTables[i].Initialize()
	}
}

func (s *Synth) Play(input chan<- float32) {
	for {
		var y float64

		for i := range s.WaveTables {
			w := s.WaveTables[i]
			y += w.SignalFunc(s.Phase) * s.Gain
		}

		s.Phase += s.step
		y /= float64(len(s.WaveTables))
		input <- float32(y)
	}
}

func (s *Synth) FadeOut() {
	sampleRate := config.Instance().SampleRate
	for s.Gain > 0 {
		s.Gain -= 0.01
		time.Sleep(time.Second / time.Duration(sampleRate))
	}
}

func (s *Synth) FadeIn() {
	sampleRate := config.Instance().SampleRate
	for s.Gain < s.gainMemory {
		s.Gain += 0.01
		time.Sleep(time.Second / time.Duration(sampleRate))
	}
}
