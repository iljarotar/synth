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

func newIntegralFunc(oscType OscillatorType) SignalFunc {
	switch oscType {
	case Sine:
		return IntegralSineSignalFunc()
	case Square:
		return IntegralSquareSignalFunc()
	case Sawtooth:
		return IntegralSawtoothSignalFunc()
	case Triangle:
		return SineSignalFunc()
	case ReverseSawtooth:
		return IntegralReverseSawtoothSignalFunc()
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

func IntegralSineSignalFunc() SignalFunc {
	integral := func(x float64) float64 {
		return -math.Cos(x)
	}

	return integral
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

func IntegralSquareSignalFunc() SignalFunc {
	integral := func(x float64) float64 {
		return math.Acos(math.Cos(x))
	}

	return integral
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

func IntegralSawtoothSignalFunc() SignalFunc {
	integral := func(x float64) float64 {
		return -math.Pi / 2 * math.Sqrt(1-math.Pow(math.Sin(x/2), 2))
	}

	return integral
}

func ReverseSawtoothSignalFunc() SignalFunc {
	sawtooth := func(x float64) float64 {
		return 1 - 2*(x/(2*math.Pi)-math.Floor(1/2+x/(2*math.Pi))) - 1
	}

	return sawtooth
}

func IntegralReverseSawtoothSignalFunc() SignalFunc {
	integral := func(x float64) float64 {
		return math.Pi / 2 * math.Sqrt(1-math.Pow(math.Sin(x/2), 2))
	}

	return integral
}
