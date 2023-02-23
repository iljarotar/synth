package utils

import "math"

// normalizes x from [-1;1] to [min;max]
func InverseNormalize(x, min, max float64) float64 {
	return Percentage(x, -1, 1)*(max-min) + min
}

// normalizes x from [min;max] to [-1;1]
func Normalize(x, min, max float64) float64 {
	return Percentage(x, min, max)*2 - 1
}

// transposes x from [min;max] to [0;1]
func Percentage(x, min, max float64) float64 {
	return (x - min) / (max - min)
}

// transposes x from [oldMin; oldMax] to [newMin; newMax]
func Transpose(x, oldMin, oldMax, newMin, newMax float64) float64 {
	return InverseNormalize(Normalize(x, oldMin, oldMax), newMin, newMax)
}

// limits x to [min;max]
func Limit(x, min, max float64) float64 {
	var y float64

	if x > min {
		y = math.Min(x, max)
	} else {
		y = min
	}

	return y
}
