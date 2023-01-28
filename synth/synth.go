package synth

import (
	w "github.com/iljarotar/synth/wavetable"
)

type Synth struct {
	Gain      float64     `yaml:"gain"`
	WaveTable w.WaveTable `yaml:"wavetable"`
}

func (s *Synth) Initialize() {
	if s.Gain == 0 {
		s.Gain = 1
	} else {
	}

	s.WaveTable.Initialize()
}

func (s *Synth) Play(input chan<- float32) {
	for {
		input <- float32(s.WaveTable.SignalFunc(s.WaveTable.Phase) * s.Gain)
		s.WaveTable.Phase += s.WaveTable.Step
	}
}

func (s *Synth) SetWaveTable(waveTable w.WaveTable) {
	s.WaveTable = waveTable
}
