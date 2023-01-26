package oscillator

import (
	"math"
	"math/rand"
	"time"
)

type SignalFunc func(x float64) float64

func NewFunc(waveType OscillatorType) SignalFunc {
	switch waveType {
	case Sine:
		return SineFunc()
	case Noise:
		return NoiseFunc()
	default:
		return NoFunc()
	}
}

func NoFunc() SignalFunc {
	f := func(x float64) float64 {
		return 0
	}
	return f
}

func SineFunc() SignalFunc {
	sine := func(x float64) float64 {
		return math.Sin(2 * math.Pi * x)
	}

	return sine
}

func NoiseFunc() SignalFunc {
	rand.Seed(time.Now().Unix())
	noise := func(x float64) float64 {
		y := rand.Float64()*2 - 1
		return y
	}

	return noise
}
