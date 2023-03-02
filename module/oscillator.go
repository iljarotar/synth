package module

import (
	"github.com/iljarotar/synth/utils"
)

type OscillatorType string

func (t OscillatorType) String() string {
	return string(t)
}

const (
	Sawtooth         OscillatorType = "Sawtooth"
	InvertedSawtooth OscillatorType = "InvertedSawtooth"
	Sine             OscillatorType = "Sine"
	Square           OscillatorType = "Square"
	Triangle         OscillatorType = "Triangle"
	Noise            OscillatorType = "Noise"
)

type limits struct {
	high, low float64
}

var (
	ampLimits      limits = limits{low: 0, high: 1}
	ampModLimits   limits = limits{low: 0, high: 1}
	panLimits      limits = limits{low: -1, high: 1}
	panModLimits   limits = limits{low: 0, high: 1}
	phaseLimits    limits = limits{low: -1, high: 1}
	phaseModLimits limits = limits{low: 0, high: 20000} // upper limit is arbitrary
)

type output struct {
	Mono, Left, Right float64
}

type Oscillators map[string]*Oscillator

type Oscillator struct {
	Name    string         `yaml:"name"`
	Type    OscillatorType `yaml:"type"`
	Freq    []float64      `yaml:"freq"`
	Amp     Param          `yaml:"amp"`
	Phase   Param          `yaml:"phase"`
	Filters []string       `yaml:"filters"`
	Pan     Param          `yaml:"pan"`
	signal  SignalFunc
	Current output
	pan     float64
}

func (o *Oscillator) Initialize() {
	o.signal = NewSignalFunc(o.Type)
	o.limit()

	var y float64
	for _, f := range o.Freq {
		y += o.partial(f, o.Phase.Val, o.Amp.Val, make(Filters))
	}

	if l := len(o.Freq); l > 0 {
		y /= float64(l)
	}

	o.Current = o.stereo(y)
}

func (o *Oscillator) Next(oscMap Oscillators, filtersMap Filters, phase float64) {
	o.pan = utils.Limit(o.Pan.Val+modulate(o.Pan.Mod, oscMap)*o.Pan.ModAmp, panLimits.low, panLimits.high)
	amp := utils.Limit(o.Amp.Val+modulate(o.Amp.Mod, oscMap)*o.Amp.ModAmp, ampLimits.low, ampLimits.high)

	if o.Type == Noise {
		o.Current = o.stereo((o.signal(0) * amp)) // noise doesn't care about phase
		return
	}

	// phase shift should not be limitted
	shift := o.Phase.Val + modulate(o.Phase.Mod, oscMap)*o.Phase.ModAmp

	var y float64
	for _, f := range o.Freq {
		l := 1 / f
		s := shift * l
		y += o.partial(f, phase+s, amp, filtersMap)
	}

	if l := len(o.Freq); l > 0 {
		y /= float64(l)
	}

	o.Current = o.stereo(y)
}

func (o *Oscillator) limit() {
	o.Amp.ModAmp = utils.Limit(o.Amp.ModAmp, ampModLimits.low, ampModLimits.high)
	o.Amp.Val = utils.Limit(o.Amp.Val, ampLimits.low, ampLimits.high)

	o.Phase.ModAmp = utils.Limit(o.Phase.ModAmp, phaseModLimits.low, phaseModLimits.high)
	o.Phase.Val = utils.Limit(o.Phase.Val, phaseLimits.low, phaseLimits.high)

	o.Pan.ModAmp = utils.Limit(o.Pan.ModAmp, panModLimits.low, panModLimits.high)
	o.Pan.Val = utils.Limit(o.Pan.Val, panLimits.low, panLimits.high)
	o.pan = o.Pan.Val

	for _, f := range o.Freq {
		f = utils.Limit(f, 0, 20000)
	}
}

func modulate(modulators []string, oscMap Oscillators) float64 {
	var y float64

	for _, m := range modulators {
		mod, ok := oscMap[m]
		if ok {
			y += mod.Current.Mono
		}
	}

	return y
}

func (o *Oscillator) applyFilters(filtersMap Filters, freq float64) float64 {
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

func (o *Oscillator) partial(freq, phase, amp float64, filtersMap Filters) float64 {
	a := amp

	if len(o.Filters) > 0 {
		a *= o.applyFilters(filtersMap, freq)
	}

	return o.signal(freq*phase) * a
}

func (o *Oscillator) stereo(x float64) output {
	out := output{}
	pan := utils.Percentage(o.pan, -1, 1)
	out.Mono = x
	out.Right = x * pan
	out.Left = x * (1 - pan)

	return out
}
