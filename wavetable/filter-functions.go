package wavetable

type FilterFunc func(freq, cutoff, ramp float64) (amp float64)

func NewFilterFunc(filterType FilterType) FilterFunc {
	switch filterType {
	case Lowpass:
		return LowpassFilterFunc()
	default:
		return NoFilterFunc()
	}
}

func NoFilterFunc() FilterFunc {
	f := func(freq, cutoff, ramp float64) (amp float64) {
		return freq
	}
	return f
}

func LowpassFilterFunc() FilterFunc {
	f := func(freq, cutoff, ramp float64) (amp float64) {
		if freq < cutoff {
			return 1
		}
		return 0 // return linear ramp instead
	}

	return f
}
