package module

import "github.com/iljarotar/synth/utils"

var (
	cutoffLimits limits = limits{low: 0, high: 20000}
	volumeLimits limits = limits{low: 0, high: 1}
)

type Filters map[string]*Filter

type Filter struct {
	Name           string  `yaml:"name"`
	Low            Param   `yaml:"low"`
	High           Param   `yaml:"high"`
	Volume         Param   `yaml:"vol"`
	Ramp           float64 `yaml:"ramp"`
	Func           FilterFunc
	low, high, vol float64
}

func (f *Filter) Initialize() {
	f.Func = NewFilterFunc(f.Ramp)

	f.Volume.Val = utils.Limit(f.Volume.Val, 0, 1)
	f.Low.Val = utils.Limit(f.Low.Val, 0, 20000)
	f.High.Val = utils.Limit(f.High.Val, 0, 20000)

	f.vol = f.Volume.Val
	f.low = f.Low.Val
	f.high = f.High.Val
}

func (f *Filter) Next(oscMap Oscillators) {
	f.low = utils.Limit(f.Low.Val+modulate(f.Low.Mod, oscMap)*f.Low.ModAmp, cutoffLimits.low, cutoffLimits.high)
	f.high = utils.Limit(f.High.Val+modulate(f.High.Mod, oscMap)*f.High.ModAmp, cutoffLimits.low, cutoffLimits.high)
	f.vol = utils.Limit(f.Volume.Val+modulate(f.Volume.Mod, oscMap)*f.Volume.ModAmp, volumeLimits.low, volumeLimits.high)
}

func (f *Filter) Apply(freq float64) float64 {
	return f.Func(freq, f.low, f.high, f.vol)
}
