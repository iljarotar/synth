package module

type FilterType string

func (t FilterType) String() string {
	return string(t)
}

const (
	Lowpass  FilterType = "Lowpass"
	Highpass FilterType = "Highpass"
)

type Filters map[string]*Filter

type Filter struct {
	Name          string     `yaml:"name"`
	Type          FilterType `yaml:"type"`
	Cutoff        Param      `yaml:"cutoff"`
	Ramp          float64    `yaml:"ramp"`
	Func          FilterFunc
	currentCutoff float64
}

func (f *Filter) Initialize() {
	f.Func = NewFilterFunc(f.Type)
	f.currentCutoff = f.Cutoff.Val
}

func (f *Filter) UpdateCutoff(x ...float64) {
	f.currentCutoff = f.Cutoff.Val
	for i := range x {
		f.currentCutoff += x[i]
	}
}

func (f *Filter) Apply(freq float64) float64 {
	return f.Func(freq, f.currentCutoff, f.Ramp)
}
