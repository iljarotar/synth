package synth

import (
	"time"

	"github.com/iljarotar/synth/config"
	w "github.com/iljarotar/synth/wavetable"
)

type Synth struct {
	Gain, gainMemory float64     `yaml:"gain"`
	WaveTable        w.WaveTable `yaml:"wavetable"`
}

func (s *Synth) Initialize() {
	if s.Gain == 0 {
		s.Gain = 1
	}
	s.gainMemory = s.Gain
	s.Gain = 0 // start muted
	s.WaveTable.Initialize()
}

func (s *Synth) Play(input chan<- float32) {
	for {
		input <- float32(s.WaveTable.SignalFunc(s.WaveTable.Phase) * s.Gain)
		s.WaveTable.Phase += s.WaveTable.Step
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
