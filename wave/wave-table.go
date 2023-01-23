package wave

import (
	"math"
	"math/rand"
	"time"

	"github.com/iljarotar/synth/config"
)

type SignalFunc func(x ...float64) float64

type WaveTable struct {
	step, phase float64
	SignalFunc  SignalFunc
}

func NewWaveTable(functions []SignalFunc) WaveTable {
	signalFunc := func(x ...float64) float64 {
		var y float64
		for i := range functions {
			functions[i](x...)
		}
		return y
	}

	c := config.Instance()
	w := WaveTable{SignalFunc: signalFunc, phase: 0, step: 1 / c.SampleRate}
	return w
}

func (w *WaveTable) Process(out []float32) {
	for i := range out {
		out[i] = float32(w.SignalFunc(w.phase))
		_, w.phase = math.Modf(w.phase + w.step)
	}
}

func SineFunc(freq float64) SignalFunc {
	sine := func(x ...float64) float64 {
		return math.Sin(2 * math.Pi * freq * x[0])
	}

	return sine
}

func NoiseFunc() SignalFunc {
	rand.Seed(time.Now().Unix())
	noise := func(x ...float64) float64 {
		return rand.Float64()
	}

	return noise
}
