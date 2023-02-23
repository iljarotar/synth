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

func (f *Filter) Next(oscMap Oscillators) {
	f.low = modulate(f.Low.Val, f.Low.Mod, oscMap)
	f.high = modulate(f.High.Val, f.High.Mod, oscMap)
	f.vol = modulate(f.Volume.Val, f.Volume.Mod, oscMap)
}

func (f *Filter) Apply(freq float64) float64 {
	return f.Func(freq, f.low, f.high, f.vol)
}
