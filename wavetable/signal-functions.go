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
	case Sawtooth:
		return SawtoothSignalFunc()
	case Triangle:
		return TriangleSignalFunc()
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
		n, _ := math.Modf(x / math.Pi)
		arg := (2*math.Pi*x - float64(n)*math.Pi)
		y := math.Sin(2 * math.Pi * x)

		if y > 0 {
			return 2/(1+math.Exp(-5*arg)) - 1
		}
		return 1 - 2/(1+math.Exp(-5*arg))
	}

	return square
}

func TriangleSignalFunc() SignalFunc {
	triangle := func(x float64) float64 {
		return 2 / math.Pi * math.Asin(math.Sin(2*math.Pi*x))
	}

	return triangle
}

func SawtoothSignalFunc() SignalFunc {
	sawtooth := func(x float64) float64 {
		return 2*(x-math.Floor(1/2+x)) - 1
	}

	return sawtooth
}

func NoiseSignalFunc() SignalFunc {
	rand.Seed(time.Now().Unix())
	noise := func(x float64) float64 {
		y := rand.Float64()*2 - 1
		return y
	}

	return noise
}
