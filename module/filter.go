package module

import (
	"math"

	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/utils"
)

type FilterType string

const (
	FilterTypeLowPass  FilterType = "LowPass"
	FilterTypeHighPass FilterType = "HighPass"
)

type Filter struct {
	Type                   FilterType `yaml:"type"`
	Cutoff                 Param      `yaml:"cutoff"`
	DBGain                 float64    `yaml:"db-gain"`
	Slope                  float64    `yaml:"slope"`
	a0, a1, a2, b0, b1, b2 float64
}

func (f *Filter) Initialize() {
	f.limitParams()
	f.calculateCoeffs(f.Cutoff.Val)
}

func (f *Filter) Tap(x2, x1, x0, y1, y0 float64) (y2 float64) {
	y2 = (f.b0/f.a0)*x2 + (f.b1/f.a0)*x1 + (f.b2/f.a0)*x0 - (f.a1/f.a0)*y1 - (f.a2/f.a0)*y0
	return y2
}

func (f *Filter) NextCoeffs(modMap ModulesMap) {
	fc := modulate(f.Cutoff, freqLimits, modMap)
	f.calculateCoeffs(fc)
}

func (f *Filter) calculateCoeffs(fc float64) {
	amp := getAmp(f.DBGain)
	omega := getOmega(fc)
	alpha := getAlpha(omega, amp, f.Slope)

	if f.Type == FilterTypeLowPass {
		f.calculateLowPassCoeffs(omega, alpha)
	} else if f.Type == FilterTypeHighPass {
		f.calculateHighPassCoeffs(omega, alpha)
	}
}

func (f *Filter) calculateLowPassCoeffs(omega, alpha float64) {
	f.b1 = 1 - math.Cos(omega)
	f.b0 = f.b1 / 2
	f.b2 = f.b0
	f.a0 = 1 + alpha
	f.a1 = -2 * math.Cos(omega)
	f.a2 = 1 - alpha
}

func (f *Filter) calculateHighPassCoeffs(omega, alpha float64) {
	f.b0 = (1 + math.Cos(omega)) / 2
	f.b1 = -(1 + math.Cos(omega))
	f.b2 = f.b0
	f.a0 = 1 + alpha
	f.a1 = -2 * math.Cos(omega)
	f.a2 = 1 - alpha
}

func getAmp(dbGain float64) float64 {
	return math.Pow(10, dbGain/40)
}

func getOmega(fc float64) float64 {
	return 2 * math.Pi * (fc / config.Config.SampleRate)
}

func getAlpha(omega, amp, slope float64) float64 {
	rootArg := (amp+1/amp)*(1/slope-slope) + 2
	root := math.Sqrt(rootArg)
	factor := math.Sin(omega) / 2
	return factor * root
}

func (f *Filter) limitParams() {
	f.Cutoff.Val = utils.Limit(f.Cutoff.Val, freqLimits.min, freqLimits.max)
	f.Cutoff.ModAmp = utils.Limit(f.Cutoff.ModAmp, freqLimits.min, freqLimits.max)
}
