package synth

import (
	o "github.com/iljarotar/synth/oscillator"
)

type Synth struct {
	Gain      float64     `yaml:"gain"`
	WaveTable o.WaveTable `yaml:"wavetable"`
}

func (s *Synth) Initialize() {
	if s.Gain == 0 {
		s.Gain = 1
	} else {
		s.Gain /= 100 // gain is given in percent
	}

	s.WaveTable.Initialize()
}

func (s *Synth) Play(input chan<- float32, play *bool) {
	for *play {
		input <- float32(s.WaveTable.SignalFunc(s.WaveTable.Phase) * s.Gain)
		s.WaveTable.Phase += s.WaveTable.Step
	}
	close(input)
}

func (s *Synth) SetWaveTable(waveTable o.WaveTable) {
	s.WaveTable = waveTable
}
