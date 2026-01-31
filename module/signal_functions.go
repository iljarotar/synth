package module

import (
	"fmt"
	"math"
)

type SignalFunc func(x float64) float64

func newSignalFunc(oscType oscillatorType) (SignalFunc, error) {
	switch oscType {
	case oscillatorTypeSine:
		return SineSignalFunc(), nil
	case oscillatorTypeSquare:
		return SquareSignalFunc(), nil
	case oscillatorTypeSawtooth:
		return SawtoothSignalFunc(), nil
	case oscillatorTypeTriangle:
		return TriangleSignalFunc(), nil
	case oscillatorTypeReverseSawtooth:
		return ReverseSawtoothSignalFunc(), nil
	default:
		return NoSignalFunc(), fmt.Errorf("unknown oscillator type %s", oscType)
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
		return 2 * (x/(2*math.Pi) - math.Floor(1.0/2+x/(2*math.Pi)))
	}

	return sawtooth
}

func ReverseSawtoothSignalFunc() SignalFunc {
	sawtooth := func(x float64) float64 {
		return 1 - 2*(x/(2*math.Pi)-math.Floor(1.0/2+x/(2*math.Pi))) - 1
	}

	return sawtooth
}
