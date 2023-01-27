package filter

type FilterType string

func (t FilterType) String() string {
	return string(t)
}

const (
	Lowpass  FilterType = "Lowpass"
	Highpass FilterType = "Highpass"
	Bandpass FilterType = "Bandpass"
)

type Filter struct {
	Type       FilterType `yaml:"type"`
	Cutoff     float64    `yaml:"cutoff"`
	Freq       float64    `yaml:"freq"`
	Range      float64    `yaml:"range"`
	filterFunc FilterFunc
}

func (f *Filter) Initialize() {
	f.filterFunc = NewFunc(f.Type)
}

func (f *Filter) Apply(x float64) float64 {
	return x
}
