package wave

import (
	"math"
)

type WaveTable struct {
	step, phase, SampleRate float64
	SignalFunc              func(x float64) float64
}

func (w *WaveTable) Process(out []float32) {
	for i := range out {
		out[i] = float32(w.SignalFunc(w.phase))
		_, w.phase = math.Modf(w.phase + w.step)
	}
}

func Sine(freq, sampleRate float64) *WaveTable {
	sine := func(x float64) float64 {
		return math.Sin(x * 2 * math.Pi * freq)
	}
	w := &WaveTable{SignalFunc: sine, phase: 0, step: 1 / sampleRate, SampleRate: sampleRate}
	return w
}

func Custom(sampleRate float64, functions []func(float64) float64) *WaveTable {
	signalFunc := func(x float64) float64 {
		var y float64
		for i := range functions {
			y += functions[i](x)
		}
		return y
	}

	w := &WaveTable{SignalFunc: signalFunc, phase: 0, step: 1 / sampleRate, SampleRate: sampleRate}
	return w
}
