package wave

import (
	"math"
	"math/rand"
	"time"

	"github.com/iljarotar/synth/config"
)

type SignalFunc func(float64) float64
type WaveFunc func(float64) float64
type NoiseFunc func() float64

type WaveTable struct {
	step, phase float64
	SignalFunc  SignalFunc
}

func NewWaveTable(waves []WaveFunc, noises []NoiseFunc) WaveTable {
	signalFunc := func(x float64) float64 {
		var y float64
		for i := range waves {
			y += waves[i](x)
		}

		for i := range noises {
			y += noises[i]()
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

func SineFunc(freq float64) WaveFunc {
	sine := func(x float64) float64 {
		return math.Sin(x * 2 * math.Pi * freq)
	}

	return sine
}

func Noise() NoiseFunc {
	rand.Seed(time.Now().Unix())
	noise := func() float64 {
		return rand.Float64()
	}

	return noise
}
