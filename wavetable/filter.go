package wavetable

type FilterType string

func (t FilterType) String() string {
	return string(t)
}

const (
	Lowpass  FilterType = "Lowpass"
	Highpass FilterType = "Highpass"
)

type Filter struct {
	Type       FilterType `yaml:"type"`
	Cutoff     float64    `yaml:"cutoff"`
	Ramp       float64    `yaml:"ramp"`
	CutoffMod  *WaveTable `yaml:"cutoff-mod"`
	filterFunc FilterFunc
}

func (f *Filter) Initialize() {
	f.filterFunc = NewFilterFunc(f.Type)

	if f.CutoffMod != nil {
		f.CutoffMod.Initialize()
	}
}

func (f *Filter) Apply(freq, x float64) float64 {
	cutoff := f.Cutoff

	if f.CutoffMod != nil {
		cutoff += f.CutoffMod.SignalFunc(x)
	}

	return f.filterFunc(freq, cutoff, f.Ramp)
}
