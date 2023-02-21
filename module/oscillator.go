package module

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

type Oscillators map[string]*Oscillator

type Oscillator struct {
	Name    string         `yaml:"name"`
	Type    OscillatorType `yaml:"type"`
	Freq    float64        `yaml:"freq"`
	Amp     Param          `yaml:"amp"`
	Phase   Param          `yaml:"phase"`
	Filters []string       `yaml:"filters"`
	Signal  SignalFunc
	Current float64
}

func (o *Oscillator) Initialize() {
	o.Signal = NewSignalFunc(o.Type)
	o.Current = o.Signal(o.Phase.Val) * o.Amp.Val
}

func (o *Oscillator) NextValue(oscMap Oscillators, filtersMap Filters, phase float64) {
	amp := o.getAmp(oscMap)
	shift := o.getPhase(oscMap)

	if len(o.Filters) > 0 {
		o.applyFilters(filtersMap, &amp)
	}

	o.Current = o.Signal(o.Freq*(phase+shift)) * amp
}

func (o *Oscillator) getAmp(oscMap Oscillators) float64 {
	amp := o.Amp.Val

	for i := range o.Amp.Mod {
		mod, ok := oscMap[o.Amp.Mod[i]]
		if ok {
			amp += mod.Current
		}
	}

	return amp
}

func (o *Oscillator) getPhase(oscMap Oscillators) float64 {
	phase := o.Phase.Val

	for j := range o.Phase.Mod {
		mod, ok := oscMap[o.Phase.Mod[j]]
		if ok {
			phase += mod.Current
		}
	}

	return phase
}

func (o *Oscillator) applyFilters(filtersMap Filters, amp *float64) {
	var max float64

	for i := range o.Filters {
		f, ok := filtersMap[o.Filters[i]]

		if ok {
			val := f.Apply(o.Freq)

			if val > max {
				max = val
			}
		}
	}

	*amp *= max
}
