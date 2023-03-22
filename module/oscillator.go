package module

import (
	"math"

	"github.com/iljarotar/synth/utils"
)

type OscillatorType string

func (t OscillatorType) String() string {
	return string(t)
}

const (
	Sawtooth        OscillatorType = "Sawtooth"
	ReverseSawtooth OscillatorType = "ReverseSawtooth"
	Sine            OscillatorType = "Sine"
	Square          OscillatorType = "Square"
	Triangle        OscillatorType = "Triangle"
)

type OscillatorsMap map[string]*Oscillator

type Oscillator struct {
	Name     string         `yaml:"name"`
	Type     OscillatorType `yaml:"type"`
	Freq     Param          `yaml:"freq"`
	Amp      Param          `yaml:"amp"`
	Phase    float64        `yaml:"phase"`
	Filters  []string       `yaml:"filters"`
	Pan      Param          `yaml:"pan"`
	signal   SignalFunc
	Integral SignalFunc
	Current  output
	pan, Phi float64
}

func (o *Oscillator) Initialize() {
	o.signal = newSignalFunc(o.Type)
	o.Integral = newIntegralFunc(o.Type)
	o.limitParams()
	o.calculateCurrentValue(o.Amp.Val, 0, 0, make(FiltersMap))
}

func (o *Oscillator) Next(oscMap OscillatorsMap, filtersMap FiltersMap, x float64) {
	o.pan = utils.Limit(o.Pan.Val+modulate(o.Pan.Mod, oscMap)*o.Pan.ModAmp, panLimits.low, panLimits.high)
	amp := utils.Limit(o.Amp.Val+modulate(o.Amp.Mod, oscMap)*o.Amp.ModAmp, ampLimits.low, ampLimits.high)
	fm := o.fm(oscMap)
	o.calculateCurrentValue(amp, x, fm, filtersMap)
}

func (o *Oscillator) fm(oscMap OscillatorsMap) float64 {
	var y float64

	for _, osc := range o.Freq.Mod {
		mod, ok := oscMap[osc]
		if ok {
			y += mod.Integral(mod.Phi) / mod.Freq.Val
		}
	}

	return y * o.Freq.ModAmp
}

func (o *Oscillator) calculateCurrentValue(amp, x, fm float64, filtersMap FiltersMap) {
	if len(o.Filters) > 0 {
		amp *= o.applyFilters(filtersMap, o.Freq.Val)
	}

	shift := o.Phase / o.Freq.Val // shift is a fraction of one period
	o.Phi = 2*math.Pi*o.Freq.Val*(x+shift) + fm
	y := o.signal(o.Phi) * amp

	o.Current = stereo(y, o.pan)
}

func (o *Oscillator) limitParams() {
	o.Amp.ModAmp = utils.Limit(o.Amp.ModAmp, modLimits.low, modLimits.high)
	o.Amp.Val = utils.Limit(o.Amp.Val, ampLimits.low, ampLimits.high)
	o.Phase = utils.Limit(o.Phase, phaseLimits.low, phaseLimits.high)

	o.Pan.ModAmp = utils.Limit(o.Pan.ModAmp, modLimits.low, modLimits.high)
	o.Pan.Val = utils.Limit(o.Pan.Val, panLimits.low, panLimits.high)
	o.pan = o.Pan.Val

	o.Freq.Val = utils.Limit(o.Freq.Val, freqLimits.low, freqLimits.high)
	o.Freq.ModAmp = utils.Limit(o.Freq.ModAmp, freqLimits.low, freqLimits.high)
}

func (o *Oscillator) applyFilters(filtersMap FiltersMap, freq float64) float64 {
	var max float64

	for _, f := range o.Filters {
		filter, ok := filtersMap[f]

		if ok {
			val := filter.Apply(freq)
			if val > max {
				max = val
			}
		}
	}

	return max
}
