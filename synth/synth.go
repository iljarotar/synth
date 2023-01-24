package synth

import (
	"github.com/iljarotar/synth/wave"
)

type Synth struct {
	WaveTable wave.WaveTable `yaml:"wavetable"`
	Playing   *bool
}

func (s *Synth) Initialize() {
	s.WaveTable.CreateSignalFunction()
	s.Playing = new(bool)
}

func (s *Synth) Play(input chan<- float32) {
	for *s.Playing {
		input <- float32(s.WaveTable.SignalFunc(s.WaveTable.Phase))
		s.WaveTable.Phase += s.WaveTable.Step
	}
	close(input)
}

func (s *Synth) SetWaveTable(waveTable wave.WaveTable) {
	s.WaveTable = waveTable
}
