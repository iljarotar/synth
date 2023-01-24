package synth

import (
	"github.com/iljarotar/synth/wave"
)

type Synth struct {
	waveTable wave.WaveTable
	Playing   *bool
}

func NewSynth(waveTable wave.WaveTable) *Synth {
	return &Synth{waveTable: waveTable, Playing: new(bool)}
}

func (s *Synth) Play(input chan<- float32) {
	for *s.Playing {
		input <- float32(s.waveTable.SignalFunc(s.waveTable.Phase))
		s.waveTable.Phase += s.waveTable.Step
	}
	close(input)
}
