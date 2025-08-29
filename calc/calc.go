package calc

import "math"

type Range struct {
	Min, Max float64
}

// Percentage returns x's position within Range in percent
func Percentage(x float64, r Range) float64 {
	return (x - r.Min) / (r.Max - r.Min)
}

// Limit limits x to fix into Range
func Limit(x float64, r Range) float64 {
	return math.Max(r.Min, math.Min(r.Max, x))
}

// Transppose projects x from one range to another maintaining equal proportions
func Transpose(x float64, from, to Range) float64 {
	p := Percentage(x, from)
	return p*(to.Max-to.Min) + to.Min
}
