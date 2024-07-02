package module

import (
	"math"

	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/utils"
)

const (
	maxOrder = 1000
	pi       = math.Pi
)

type Filter struct {
	Order           int     `yaml:"order"`
	LowCutoff       float64 `yaml:"low-cutoff"`
	HighCutoff      float64 `yaml:"high-cutoff"`
	weights, buffer []float64
}

func (f *Filter) Initialize() {
	f.Order = int(utils.Limit(float64(f.Order), 0, maxOrder))

	if f.LowCutoff == 0 && f.HighCutoff == 0 {
		f.weights = make([]float64, 0) // disable filter
		f.Order = 0
		return
	}

	f.buffer = make([]float64, f.Order)

	if f.LowCutoff == 0 {
		f.weights = getLowpassCoefficients(f.Order, f.HighCutoff)
		return
	}

	if f.HighCutoff == 0 {
		f.weights = getHighpassCoefficients(f.Order, f.LowCutoff)
		return
	}

	f.weights = getBandpassCoefficients(f.Order, f.LowCutoff, f.HighCutoff)
}

func (f *Filter) Tap(x float64) float64 {
	if f.Order == 0 {
		return x
	}

	f.buffer = append([]float64{x}, f.buffer[:len(f.buffer)-1]...)
	y := f.weights[0]

	for n := range f.buffer {
		y += f.weights[n+1] * f.buffer[n]
	}

	return y
}

func getLowpassCoefficients(order int, cutoff float64) []float64 {
	b := make([]float64, order+1)
	cutoff /= config.Config.SampleRate
	w := 2 * pi * cutoff

	for n := range b {
		k := float64(n - order/2)
		b[n] = (w / pi) * sinc(w*k) * hamming(n, order)
	}

	b = utils.Normalize(b, -0.22, 1) // approximately sinc range

	return b
}

func getHighpassCoefficients(order int, cutoff float64) []float64 {
	b := make([]float64, order+1)
	cutoff /= config.Config.SampleRate
	w := 2 * pi * cutoff

	for n := range b {
		k := float64(n - order/2)
		b[n] = delta(int(k)) - (w/pi)*sinc(w*k)*hamming(n, order)
	}

	b = utils.Normalize(b, -0.22, 1) // approximately sinc range

	return b
}

func getBandpassCoefficients(order int, lowCutoff, highCutoff float64) []float64 {
	b := make([]float64, order+1)
	lowCutoff /= config.Config.SampleRate
	highCutoff /= config.Config.SampleRate
	wl := 2 * pi * lowCutoff
	wh := 2 * pi * highCutoff

	for n := range b {
		k := float64(n - order/2)
		b[n] = (1 / pi) * (wh*sinc(wh*k) - wl*sinc(wl*k)) * hamming(n, order)
	}

	b = utils.Normalize(b, -0.22, 1) // approximately sinc range

	return b
}

func sinc(x float64) float64 {
	if x == 0 {
		return 1
	}
	return math.Sin(x) / x
}

func hamming(n, order int) float64 {
	hamming := 0.54 - 0.46*math.Cos(2*pi*float64(n)/float64(order))
	return hamming
}

func delta(n int) float64 {
	if n == 0 {
		return 1
	}
	return 0
}
