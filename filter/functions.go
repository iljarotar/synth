package filter

type FilterFunc func(float64) float64

func NewFunc(filterType FilterType) FilterFunc {
	return NoFunc()
}

func NoFunc() FilterFunc {
	f := func(x float64) float64 {
		return x
	}
	return f
}
