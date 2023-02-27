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
