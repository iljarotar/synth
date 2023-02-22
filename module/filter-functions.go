package module

import (
	"math"
)

type FilterFunc func(freq, low, high, vol float64) (amp float64)

func NewFilterFunc(ramp float64) FilterFunc {
	f := func(freq, low, high, vol float64) (amp float64) {
		if freq >= low && freq <= high {
			return vol
		}

		if freq < low {
			m := vol / ramp
			t := vol * (ramp - low) / ramp
			y := m*freq + t // linear ramp

			return math.Max(y, 0)
		}

		m := -vol / ramp
		t := vol * (high + ramp) / ramp
		y := m*freq + t // linear ramp

		return math.Max(y, 0)
	}

	return f
}
