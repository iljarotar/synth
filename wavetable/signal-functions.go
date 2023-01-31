package wavetable

import (
	"math"
	"math/rand"
	"time"
)

type SignalFunc func(x float64) float64

func NewSignalFunc(oscType OscillatorType) SignalFunc {
	switch oscType {
	case Sine:
		return SineSignalFunc()
	case Square:
		return SquareSignalFunc()
	case SmoothSquare:
		return SmoothSquareSignalFunc()
	case Noise:
		return NoiseSignalFunc()
	default:
		return NoSignalFunc()
	}
}

func NoSignalFunc() SignalFunc {
	f := func(x float64) float64 {
		return 0
	}
	return f
}

func SineSignalFunc() SignalFunc {
	sine := func(x float64) float64 {
		return math.Sin(2 * math.Pi * x)
	}

	return sine
}

func SquareSignalFunc() SignalFunc {
	square := func(x float64) float64 {
		y := math.Sin(2 * math.Pi * x)
		if y > 0 {
			return 1
		}
		return -1
	}

	return square
}

func SmoothSquareSignalFunc() SignalFunc {
	square := func(x float64) float64 {
		arg := 2 * math.Pi * x
		var y float64

		for i := 1; i < 12; i += 2 {
			y += math.Sin(float64(i) * arg)
		}
		return y
	}

	return square
}

func NoiseSignalFunc() SignalFunc {
	rand.Seed(time.Now().Unix())
	noise := func(x float64) float64 {
		y := rand.Float64()*2 - 1
		return y
	}

	return noise
}
