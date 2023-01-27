package filter

type FilterFunc func(freq, cutoff, ramp float64) (amp float64)

func NewFunc(filterType FilterType) FilterFunc {
	return NoFunc()
}

func NoFunc() FilterFunc {
	f := func(freq, cutoff, ramp float64) (amp float64) {
		return freq
	}
	return f
}

func LowpassFunc() FilterFunc {
	f := func(freq, cutoff, ramp float64) (amp float64) {
		if freq < cutoff {
			return freq
		}
		return 0 // return linear ramp instead
	}

	return f
}
