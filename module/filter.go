package module

type Filters map[string]*Filter

type Filter struct {
	Name           string  `yaml:"name"`
	Low            Param   `yaml:"low"`
	High           Param   `yaml:"high"`
	Volume         Param   `yaml:"vol"`
	Ramp           float64 `yaml:"ramp"`
	Func           FilterFunc
	low, high, vol float64
}

func (f *Filter) Initialize() {
	f.Func = NewFilterFunc(f.Ramp)
	f.vol = f.Volume.Val
	f.low = f.Low.Val
	f.high = f.High.Val
}

func (f *Filter) Next(oscMap Oscillators, phase float64) {
	f.low = f.Low.Val
	f.high = f.High.Val
	f.vol = f.Volume.Val

	for i := range f.Low.Mod {
		osc, ok := oscMap[f.Low.Mod[i]]
		if ok {
			f.low += osc.Current.Mono
		}
	}

	for i := range f.High.Mod {
		osc, ok := oscMap[f.High.Mod[i]]
		if ok {
			f.high += osc.Current.Mono
		}
	}

	for i := range f.Volume.Mod {
		osc, ok := oscMap[f.Volume.Mod[i]]
		if ok {
			f.vol += osc.Current.Mono
		}
	}
}

func (f *Filter) Apply(freq float64) float64 {
	return f.Func(freq, f.low, f.high, f.vol)
}
