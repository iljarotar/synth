package module

import (
	"math"
)

type FilterFunc func(freq, cutoff, ramp float64) (amp float64)

func NewFilterFunc(filterType FilterType) FilterFunc {
	switch filterType {
	case Lowpass:
		return LowpassFilterFunc()
	case Highpass:
		return HighpassFilterFunc()
	default:
		return NoFilterFunc()
	}
}

func NoFilterFunc() FilterFunc {
	f := func(freq, cutoff, ramp float64) (amp float64) {
		return 1
	}
	return f
}

func LowpassFilterFunc() FilterFunc {
	f := func(freq, cutoff, ramp float64) (amp float64) {
		if freq <= cutoff {
			return 1
		}

		m := -1 / ramp
		t := (cutoff + ramp) / ramp
		y := m*freq + t // linear ramp

		return math.Max(y, 0)
	}

	return f
}

func HighpassFilterFunc() FilterFunc {
	f := func(freq, cutoff, ramp float64) (amp float64) {
		if freq >= cutoff {
			return 1
		}

		m := 1 / ramp
		t := 1 - cutoff/ramp
		y := m*freq + t // linear ramp

		return math.Max(y, 0)
	}

	return f
}
