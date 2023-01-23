package wave

import (
	"math"

	"github.com/iljarotar/synth/config"
)

type SignalFunc func(float64) float64

type WaveTable struct {
	step, phase float64
	Config      *config.Config
	SignalFunc  SignalFunc
}

func (w *WaveTable) Process(out []float32) {
	for i := range out {
		out[i] = float32(w.SignalFunc(w.phase))
		_, w.phase = math.Modf(w.phase + w.step)
	}
}

func Sine(c *config.Config, freq float64) *WaveTable {
	sine := func(x float64) float64 {
		return math.Sin(x * 2 * math.Pi * freq)
	}
	w := &WaveTable{SignalFunc: sine, phase: 0, step: 1 / c.SampleRate, Config: c}
	return w
}

func Custom(c *config.Config, functions []SignalFunc) *WaveTable {
	signalFunc := func(x float64) float64 {
		var y float64
		for i := range functions {
			y += functions[i](x)
		}
		return y
	}

	w := &WaveTable{SignalFunc: signalFunc, phase: 0, step: 1 / c.SampleRate, Config: c}
	return w
}
