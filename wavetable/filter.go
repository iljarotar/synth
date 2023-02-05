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
	Ramp       float64    `yaml:"ramp"`
	Cutoff     *Param     `yaml:"cutoff"`
	filterFunc FilterFunc
}

func (f *Filter) Initialize() {
	f.filterFunc = NewFilterFunc(f.Type)

	if f.Cutoff != nil && f.Cutoff.Modulation != nil {
		f.Cutoff.Modulation.Initialize()
	}
}

func (f *Filter) Apply(freq, x float64) float64 {
	cutoff := f.Cutoff.Value

	if f.Cutoff != nil && f.Cutoff.Modulation != nil {
		cutoff += f.Cutoff.Modulation.SignalFunc(x)
	}

	return f.filterFunc(freq, cutoff, f.Ramp)
}
