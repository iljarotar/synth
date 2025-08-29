package utils

import "math"

func Percentage(x, min, max float64) float64 {
	return (x - min) / (max - min)
}

func Limit(x, min, max float64) float64 {
	return math.Max(min, math.Min(max, x))
}
