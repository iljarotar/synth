package wave

import (
	"math"
	"math/rand"
	"time"
)

type SignalFunc func(x ...float64) float64

func SineFunc(freq int, amplitude float64) SignalFunc {
	sine := func(x ...float64) float64 {
		return math.Sin(2*math.Pi*float64(freq)*x[0]) * amplitude
	}

	return sine
}

func NoiseFunc(amplitude float64) SignalFunc {
	rand.Seed(time.Now().Unix())
	noise := func(x ...float64) float64 {
		return rand.Float64() * amplitude
	}

	return noise
}
