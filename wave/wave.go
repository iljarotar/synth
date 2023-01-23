package wave

import (
	"math"

	"github.com/iljarotar/synth/config"
)

type WaveType string

func (t WaveType) String() string {
	return string(t)
}

const (
	Sine  WaveType = "Sine"
	Noise WaveType = "Noise"
)

type WaveTable struct {
	step, phase float64
	SignalFunc  SignalFunc
}

type Wave struct {
	Type            WaveType
	Freq, Amplitude *float64
}

// TODO: implement SignalFunc initialization and normalization
// Sum all the amplitudes to get maximum possible amplitude
func NewWaveTable(waves ...Wave) WaveTable {
	signalFunc := func(x ...float64) float64 {
		return 0
	}

	c := config.Instance()
	return WaveTable{step: 1 / c.SampleRate, SignalFunc: signalFunc}
}

func (w *WaveTable) Process(out []float32) {
	for i := range out {
		out[i] = float32(w.SignalFunc(w.phase))
		_, w.phase = math.Modf(w.phase + w.step)
	}
}
