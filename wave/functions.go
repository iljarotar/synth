package wave

import (
	"math"
	"math/rand"
	"time"
)

type SignalFunc func(x ...float64) float64

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
