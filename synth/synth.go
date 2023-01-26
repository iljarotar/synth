package synth

import (
	"math"

	o "github.com/iljarotar/synth/oscillator"
)

type Synth struct {
	WaveTable o.WaveTable `yaml:"wavetable"`
}

func (s *Synth) Initialize() {
	s.WaveTable.Initialize()
}

func (s *Synth) Play(input chan<- float32, play *bool) {
	for *play {
		input <- float32(s.WaveTable.SignalFunc(s.WaveTable.Phase))
		_, s.WaveTable.Phase = math.Modf(s.WaveTable.Phase + s.WaveTable.Step)
	}
	close(input)
}

func (s *Synth) SetWaveTable(waveTable o.WaveTable) {
	s.WaveTable = waveTable
}
