package module

import "math"

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

type Output struct {
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
	Current Output
	pan     float64
}

func (o *Oscillator) Initialize() {
	o.signal = NewSignalFunc(o.Type)
	o.pan = o.Pan.Val
	var y float64

	for i := range o.Freq {
		y += o.partial(o.Freq[i], o.Phase.Val, o.Amp.Val, make(Filters))
		o.Current = o.stereo(y)
	}
}

func (o *Oscillator) Next(oscMap Oscillators, filtersMap Filters, phase float64) {
	o.pan = modulate(o.Pan.Val, o.Pan.Mod, oscMap)
	amp := modulate(o.Amp.Val, o.Amp.Mod, oscMap)

	if o.Type == Noise {
		o.Current = o.stereo((o.signal(0) * amp)) // noise doesn't care about phase
		return
	}

	shift := modulate(o.Phase.Val, o.Phase.Mod, oscMap)
	var y float64

	for i := range o.Freq {
		y += o.partial(o.Freq[i], phase+shift, amp, filtersMap)
	}

	y /= float64(len(o.Freq))
	o.Current = o.stereo(y)
}

func modulate(initial float64, modulators []string, oscMap Oscillators) float64 {
	new := initial

	for i := range modulators {
		mod, ok := oscMap[modulators[i]]
		if ok {
			new += mod.Current.Mono
		}
	}

	return new
}

func (o *Oscillator) applyFilters(filtersMap Filters, freq, amp float64) float64 {
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

	return max * amp
}

func (o *Oscillator) partial(freq, phase, amp float64, filtersMap Filters) float64 {
	a := amp

	if len(o.Filters) > 0 {
		a = o.applyFilters(filtersMap, freq, amp)
	}

	return o.signal(freq*phase) * a
}

func (o *Oscillator) stereo(x float64) Output {
	out := Output{}
	pan := transpose(o.pan)
	out.Mono = x
	out.Right = x * pan
	out.Left = x * (1 - pan)

	return out
}

// limits pan to [-1;1] and transposes to [0;1]
func transpose(pan float64) float64 {
	var t float64

	if pan > -1 {
		t = math.Min(pan, 1)
	} else {
		t = -1
	}

	return (t + 1) / 2
}
