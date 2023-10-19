package module

import (
	"math"

	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/utils"
)

type Filter struct {
	Order           int     `yaml:"order"`
	LowCutoff       float64 `yaml:"low-cutoff"`
	HighCutoff      float64 `yaml:"high-cutoff"`
	weights, buffer []float64
}

func (f *Filter) Initialize() {
	if f.Order < 0 {
		f.Order = 0
	}

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
	pi := math.Pi

	for n := range b {
		k := float64(n - order/2)
		hamming := 0.54 - 0.46*math.Cos(2*pi*float64(n)/float64(order))

		if k == 0 {
			b[n] = 2 * cutoff
		} else {
			b[n] = math.Sin(2*pi*cutoff*k) / (pi * k)
		}

		b[n] *= hamming
	}

	b = utils.Normalize(b, -0.22, 1) // approximately sinc range

	return b
}

func getHighpassCoefficients(order int, cutoff float64) []float64 {
	b := make([]float64, order+1)
	cutoff /= config.Config.SampleRate
	pi := math.Pi

	for n := range b {
		k := float64(n - order/2)
		hamming := 0.54 - 0.46*math.Cos(2*pi*float64(n)/float64(order))

		if k == 0 {
			b[n] = 1 - 2*cutoff
		} else {
			b[n] = -math.Sin(2*pi*cutoff*k) / (pi * k)
		}

		b[n] *= hamming
	}

	b = utils.Normalize(b, -0.22, 1) // approximately sinc range

	return b
}

func getBandpassCoefficients(order int, lowCutoff, highCutoff float64) []float64 {
	b := make([]float64, order+1)
	lowCutoff /= config.Config.SampleRate
	highCutoff /= config.Config.SampleRate
	pi := math.Pi

	for n := range b {
		k := float64(n - order/2)
		hamming := 0.54 - 0.46*math.Cos(2*pi*float64(n)/float64(order))

		if k == 0 {
			b[n] = 2 * (highCutoff - lowCutoff)
		} else {
			b[n] = (math.Sin(2*pi*highCutoff*k) - math.Sin(2*pi*lowCutoff*k)) / (pi * k)
		}

		b[n] *= hamming
	}

	b = utils.Normalize(b, -0.22, 1) // approximately sinc range

	return b
}
