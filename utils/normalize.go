package utils

import "math"

// transposes x from [min;max] to [0;1]
func Percentage(x, min, max float64) float64 {
	return (x - min) / (max - min)
}

// limits x to [min;max]
func Limit(x, min, max float64) float64 {
	y := x

	if y > min {
		y = math.Min(x, max)
	} else {
		y = min
	}

	return y
}

func Normalize(signal []float64, min, max float64) []float64 {
	var m, M float64

	for _, y := range signal {
		if y > M {
			M = y
			continue
		}
		if y < m {
			m = y
		}
	}

	for i, y := range signal {
		signal[i] = Percentage(y, m, M)*(max-min) + min
	}

	return signal
}
