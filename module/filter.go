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

func (f *Filter) UpdateCutoff(oscMap Oscillators, phase float64) {
	f.currentCutoff = f.Cutoff.Val

	for i := range f.Cutoff.Mod {
		osc, ok := oscMap[f.Cutoff.Mod[i]]
		if ok {
			f.currentCutoff += osc.Current
		}
	}
}

func (f *Filter) Apply(freq float64) float64 {
	return f.Func(freq, f.currentCutoff, f.Ramp)
}
