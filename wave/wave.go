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
	Type      WaveType
	Amplitude float64
	Freq      int
}

func NewWaveTable(waves ...Wave) WaveTable {
	var amp float64
	f := make([]SignalFunc, 0)

	for i := range waves {
		w := waves[i]
		amp += w.Amplitude
		f = append(f, NewFunc(w.Freq, w.Amplitude, w.Type))
	}

	signalFunc := func(x ...float64) float64 {
		var y float64
		for i := range waves {
			y += f[i](x...)
		}
		return y / amp // normalize
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
