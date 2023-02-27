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
	for i := range o.Freq {
		y += o.partial(o.Freq[i], o.Phase.Val, o.Amp.Val, make(Filters))
	}

	if len(o.Freq) > 0 {
		y /= float64(len(o.Freq))
	}

	o.Current = o.stereo(y)
}

func (o *Oscillator) Next(oscMap Oscillators, filtersMap Filters, phase float64) {
	o.pan = utils.Limit(o.Pan.Val+modulate(o.Pan.Mod, oscMap)*o.Pan.ModAmp, -1, 1)
	amp := utils.Limit(o.Amp.Val+modulate(o.Amp.Mod, oscMap)*o.Amp.ModAmp, 0, 1)

	if o.Type == Noise {
		o.Current = o.stereo((o.signal(0) * amp)) // noise doesn't care about phase
		return
	}

	// phase shift should not be limitted
	shift := o.Phase.Val + modulate(o.Phase.Mod, oscMap)*o.Phase.ModAmp

	var y float64
	for i := range o.Freq {
		f := o.Freq[i]
		l := 1 / f
		s := shift * l
		y += o.partial(f, phase+s, amp, filtersMap)
	}

	if len(o.Freq) > 0 {
		y /= float64(len(o.Freq))
	}

	o.Current = o.stereo(y)
}

func (o *Oscillator) limit() {
	o.Amp.ModAmp = utils.Limit(o.Amp.ModAmp, 0, 1)
	o.Amp.Val = utils.Limit(o.Amp.Val, 0, 1)

	// upper limit of modamp is arbitrary but should be sufficiently high because the lower the modulators frequency
	// the weaker the impact of the modulation will be
	o.Phase.ModAmp = utils.Limit(o.Phase.ModAmp, 0, 20000)
	o.Phase.Val = utils.Limit(o.Phase.Val, -1, 1)

	o.Pan.ModAmp = utils.Limit(o.Pan.ModAmp, 0, 1)
	o.Pan.Val = utils.Limit(o.Pan.Val, -1, 1)
	o.pan = o.Pan.Val

	for i := range o.Freq {
		o.Freq[i] = utils.Limit(o.Freq[i], 0, 20000)
	}
}

func modulate(modulators []string, oscMap Oscillators) float64 {
	var y float64

	for i := range modulators {
		mod, ok := oscMap[modulators[i]]
		if ok {
			y += mod.Current.Mono
		}
	}

	return y
}

func (o *Oscillator) applyFilters(filtersMap Filters, freq float64) float64 {
	var max float64

	for i := range o.Filters {
		f, ok := filtersMap[o.Filters[i]]

		if ok {
			val := f.Apply(freq)

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
