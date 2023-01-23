package wave

import (
	"math"

	"github.com/iljarotar/synth/config"
)

type SignalFunc func(float64) float64

type WaveTable struct {
	step, phase float64
	SignalFunc  SignalFunc
}

func (w *WaveTable) Process(out []float32) {
	for i := range out {
		out[i] = float32(w.SignalFunc(w.phase))
		_, w.phase = math.Modf(w.phase + w.step)
	}
}

func Sine(freq float64) *WaveTable {
	sine := func(x float64) float64 {
		return math.Sin(x * 2 * math.Pi * freq)
	}

	c := config.Instance()
	w := &WaveTable{SignalFunc: sine, phase: 0, step: 1 / c.SampleRate}
	return w
}

func Custom(functions []SignalFunc) *WaveTable {
	signalFunc := func(x float64) float64 {
		var y float64
		for i := range functions {
			y += functions[i](x)
		}
		return y
	}

	c := config.Instance()
	w := &WaveTable{SignalFunc: signalFunc, phase: 0, step: 1 / c.SampleRate}
	return w
}
