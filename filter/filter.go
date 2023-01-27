package filter

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
	filterFunc FilterFunc
}

func (f *Filter) Initialize() {
	f.filterFunc = NewFunc(f.Type)
}

func (f *Filter) Apply(x float64) float64 {
	return f.filterFunc(x, f.Cutoff, f.Ramp)
}