package wave

import (
	"math"
	"math/rand"
	"time"
)

type SignalFunc func(x ...float64) float64

func NewFunc(freq int, amplitude float64, waveType WaveType) SignalFunc {
	switch waveType {
	case Sine:
		return SineFunc(freq, amplitude)
	case Noise:
		return NoiseFunc(amplitude)
	default:
		return NoFunc()
	}
}

func NoFunc() SignalFunc {
	f := func(x ...float64) float64 {
		return 0
	}
	return f
}

func SineFunc(freq int, amplitude float64) SignalFunc {
	sine := func(x ...float64) float64 {
		return math.Sin(2*math.Pi*float64(freq)*x[0]) * amplitude
	}

	return sine
}

func NoiseFunc(amplitude float64) SignalFunc {
	rand.Seed(time.Now().Unix())
	noise := func(x ...float64) float64 {
		y := rand.Float64()*2 - 1
		return y * amplitude
	}

	return noise
}
