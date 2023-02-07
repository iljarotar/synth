package synth

import (
	"time"

	"github.com/iljarotar/synth/config"
	w "github.com/iljarotar/synth/wavetable"
)

type Synth struct {
	Volume                    float64 `yaml:"volume"`
	volumeMemory, Phase, step float64
	WaveTables                []*w.WaveTable `yaml:"wavetables"`
}

func (s *Synth) Initialize() {
	sampleRate := config.Instance.SampleRate()
	s.step = 1 / sampleRate

	if s.Volume == 0 {
		s.Volume = 1
	}
	s.volumeMemory = s.Volume
	s.Volume = 0 // start muted

	for i := range s.WaveTables {
		s.WaveTables[i].Initialize()
	}
}

func (s *Synth) Play(input chan<- float32) {
	for {
		var y float64

		for i := range s.WaveTables {
			w := s.WaveTables[i]
			y += w.SignalFunc(s.Phase) * s.Volume
		}

		s.Phase += s.step

		if len(s.WaveTables) > 0 {
			y /= float64(len(s.WaveTables))
		}
		input <- float32(y)
	}
}

func (s *Synth) FadeOut() {
	sampleRate := config.Instance.SampleRate()
	for s.Volume > 0 {
		s.Volume -= 0.005
		time.Sleep(time.Second / time.Duration(sampleRate))
	}
}

func (s *Synth) FadeIn() {
	sampleRate := config.Instance.SampleRate()
	for s.Volume < s.volumeMemory {
		s.Volume += 0.005
		time.Sleep(time.Second / time.Duration(sampleRate))
	}
}
