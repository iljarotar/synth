package utils

import "math"

func Percentage(x, min, max float64) float64 {
	return (x - min) / (max - min)
}

func Limit(x, min, max float64) float64 {
	y := x

	if y > min {
		y = math.Min(x, max)
	} else {
		y = min
	}

	return y
}
