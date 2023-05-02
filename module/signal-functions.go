package module

import (
	"math"
)

type SignalFunc func(x float64) float64

func newSignalFunc(oscType OscillatorType) SignalFunc {
	switch oscType {
	case Sine:
		return SineSignalFunc()
	case Square:
		return SquareSignalFunc()
	case Sawtooth:
		return SawtoothSignalFunc()
	case Triangle:
		return TriangleSignalFunc()
	case ReverseSawtooth:
		return ReverseSawtoothSignalFunc()
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
		return math.Sin(x)
	}

	return sine
}

func SquareSignalFunc() SignalFunc {
	square := func(x float64) float64 {
		y := math.Sin(x)
		if y > 0 {
			return 1
		}
		return -1
	}

	return square
}

func TriangleSignalFunc() SignalFunc {
	triangle := func(x float64) float64 {
		return 2 / math.Pi * math.Asin(math.Sin(x))
	}

	return triangle
}

func SawtoothSignalFunc() SignalFunc {
	sawtooth := func(x float64) float64 {
		return 2 * (x/(2*math.Pi) - math.Floor(1/2+x/(2*math.Pi)))
	}

	return sawtooth
}

func ReverseSawtoothSignalFunc() SignalFunc {
	sawtooth := func(x float64) float64 {
		return 1 - 2*(x/(2*math.Pi)-math.Floor(1/2+x/(2*math.Pi))) - 1
	}

	return sawtooth
}
