package module

import "github.com/iljarotar/synth/utils"

type Filter struct {
	LowCutoff              Param `yaml:"low-cutoff"`
	HighCutoff             Param `yaml:"high-cutoff"`
	a0, a1, a2, b0, b1, b2 float64
}

func (f *Filter) Initialize() {
	f.limitParams()
	f.calculateCoeffs(f.LowCutoff.Val, f.HighCutoff.Val)
}

func (f *Filter) Tap(x2, x1, x0, y1, y0 float64) float64 {
	return 0
}

func (f *Filter) NextCoeffs(modMap ModulesMap) {
	fl := utils.Limit(f.LowCutoff.Val+modulate(f.LowCutoff.Mod, modMap)*f.LowCutoff.ModAmp, freqLimits.min, freqLimits.max)
	fh := utils.Limit(f.HighCutoff.Val+modulate(f.HighCutoff.Mod, modMap)*f.HighCutoff.ModAmp, freqLimits.min, freqLimits.max)
	f.calculateCoeffs(fl, fh)
}

func (f *Filter) calculateCoeffs(fl, fh float64) {

}

func (f *Filter) limitParams() {
	f.LowCutoff.Val = utils.Limit(f.LowCutoff.Val, freqLimits.min, freqLimits.max)
	f.LowCutoff.ModAmp = utils.Limit(f.LowCutoff.ModAmp, freqLimits.min, freqLimits.max)
	f.HighCutoff.Val = utils.Limit(f.HighCutoff.Val, freqLimits.min, freqLimits.max)
	f.HighCutoff.ModAmp = utils.Limit(f.HighCutoff.ModAmp, freqLimits.min, freqLimits.max)
}
