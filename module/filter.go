package module

import (
	"math"

	"github.com/iljarotar/synth/utils"
)

type FilterType string

type FiltersMap map[string]*Filter

type filterInputs struct {
	x0, x1, x2, y0, y1 float64
}

type filterConfig struct {
	filterNames []string
	inputs      []filterInputs
	FiltersMap
}

const (
	dbGain = -50
	slope  = 0.99
)

type Filter struct {
	Name                   string `yaml:"name"`
	LowCutoff              Input  `yaml:"low-cutoff"`
	HighCutoff             Input  `yaml:"high-cutoff"`
	a0, a1, a2, b0, b1, b2 float64
	amp                    float64
	bypass                 bool
	sampleRate             float64
}

func (f *Filter) Initialize(sampleRate float64) {
	f.sampleRate = sampleRate
	f.limitParams()
	f.amp = getAmp(dbGain)
	f.calculateCoeffs(f.LowCutoff.Val, f.HighCutoff.Val)
}

func (f *Filter) Tap(x2, x1, x0, y1, y0 float64) (y2 float64) {
	if isUnset(f.LowCutoff, cutoffLimits) && isUnset(f.HighCutoff, cutoffLimits) {
		return x2
	}
	y2 = (f.b0/f.a0)*x2 + (f.b1/f.a0)*x1 + (f.b2/f.a0)*x0 - (f.a1/f.a0)*y1 - (f.a2/f.a0)*y0
	return y2
}

func (f *Filter) NextCoeffs(modMap ModulesMap) {
	fl := modulate(f.LowCutoff, cutoffLimits, modMap)
	fh := modulate(f.HighCutoff, cutoffLimits, modMap)
	f.calculateCoeffs(fl, fh)
}

func (f *Filter) calculateCoeffs(fl, fh float64) {
	switch {
	case isUnset(f.LowCutoff, cutoffLimits):
		f.calculateLowPassCoeffs(fh)
	case isUnset(f.HighCutoff, cutoffLimits):
		f.calculateHighPassCoeffs(fl)
	default:
		f.calculateBandPassCoeffs(fl, fh)
	}
}

func (f *Filter) calculateLowPassCoeffs(fc float64) {
	omega := getOmega(fc, f.sampleRate)
	alpha := getAlphaLPHP(omega, f.amp, slope)
	f.b1 = 1 - math.Cos(omega)
	f.b0 = f.b1 / 2
	f.b2 = f.b0
	f.a0 = 1 + alpha
	f.a1 = -2 * math.Cos(omega)
	f.a2 = 1 - alpha
}

func (f *Filter) calculateHighPassCoeffs(fc float64) {
	omega := getOmega(fc, f.sampleRate)
	alpha := getAlphaLPHP(omega, f.amp, slope)
	f.b0 = (1 + math.Cos(omega)) / 2
	f.b1 = -(1 + math.Cos(omega))
	f.b2 = f.b0
	f.a0 = 1 + alpha
	f.a1 = -2 * math.Cos(omega)
	f.a2 = 1 - alpha
}

func (f *Filter) calculateBandPassCoeffs(fl, fh float64) {
	if fl > fh {
		return
	}
	bw := math.Log2(fh / fl)
	fc := fl + (fh-fl)/2
	omega := getOmega(fc, f.sampleRate)
	alpha := getAlphaBP(omega, bw)
	f.b0 = alpha
	f.b1 = 0
	f.b2 = -alpha
	f.a0 = 1 + alpha
	f.a1 = -2 * math.Cos(omega)
	f.a2 = 1 - alpha
}

func getAmp(dbGain float64) float64 {
	return math.Pow(10, dbGain/40)
}

func getOmega(fc float64, sampleRate float64) float64 {
	return 2 * math.Pi * (fc / sampleRate)
}

func getAlphaLPHP(omega, amp, slope float64) float64 {
	rootArg := (amp+1/amp)*(1/slope-slope) + 2
	root := math.Sqrt(rootArg)
	factor := math.Sin(omega) / 2
	return factor * root
}

func getAlphaBP(omega, bandwidth float64) float64 {
	a := math.Log10(2) / 2
	b := omega / math.Sin(omega)
	sinh := math.Sinh(a * b * bandwidth)
	return math.Sin(omega) * sinh
}

func (c *filterConfig) applyFilters(x float64) (float64, []filterInputs) {
	var y2, y float64
	newInputs := c.inputs

	for i, f := range c.filterNames {
		filter, ok := c.FiltersMap[f]
		if !ok {
			continue
		}
		if len(c.inputs) != len(c.filterNames) {
			return 0, c.inputs
		}

		in := c.inputs[i]
		y2 = filter.Tap(in.x2, in.x1, in.x0, in.y1, in.y0)
		y += y2

		in.x0 = in.x1
		in.x1 = in.x2
		in.y0 = in.y1
		in.y1 = y2
		in.x2 = x
		newInputs[i] = in
	}

	if len(c.filterNames) == 0 {
		y = x
	} else {
		y /= float64(len(c.filterNames))
	}

	return y, newInputs
}

func (f *Filter) limitParams() {
	f.HighCutoff.Val = utils.Limit(f.HighCutoff.Val, cutoffLimits.min, cutoffLimits.max)
	f.HighCutoff.ModAmp = utils.Limit(f.HighCutoff.ModAmp, cutoffLimits.min, cutoffLimits.max)

	maxLowCutoff := cutoffLimits.max
	if !isUnset(f.HighCutoff, cutoffLimits) {
		maxLowCutoff = f.HighCutoff.Val
	}
	f.LowCutoff.Val = utils.Limit(f.LowCutoff.Val, cutoffLimits.min, maxLowCutoff)
	f.LowCutoff.ModAmp = utils.Limit(f.LowCutoff.ModAmp, cutoffLimits.min, cutoffLimits.max)
}

func isUnset(p Input, lim limits) bool {
	return p.Val == lim.min && len(p.Mod) == 0
}
