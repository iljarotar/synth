package module

import "github.com/iljarotar/synth/calc"

type IModule interface {
	Current() Output
}

type ModuleMap map[string]IModule

type Module struct {
	current Output
}

func (m *Module) Current() Output {
	return m.current
}

var (
	gainLimits = calc.Range{
		Min: 0,
		Max: 1,
	}
	outputLimits = calc.Range{
		Min: -1,
		Max: 1,
	}
	freqLimits = calc.Range{
		Min: 0,
		Max: 20000,
	}
	panLimits = calc.Range{
		Min: -1,
		Max: 1,
	}
	bpmLimits = calc.Range{
		Min: 0,
		Max: 2000,
	}
)

type Output struct {
	Mono, Left, Right float64
}

func modulate(x float64, rng calc.Range, mod float64) float64 {
	if mod == 0 {
		return x
	}

	transposed := calc.Transpose(x, rng, outputLimits)
	transposed += mod
	transposed = calc.Limit(transposed, outputLimits)
	return calc.Transpose(transposed, outputLimits, rng)
}

func cv(rng calc.Range, val float64) float64 {
	val = calc.Limit(val, outputLimits)
	return calc.Transpose(val, outputLimits, rng)
}

func getMono(mod IModule) float64 {
	if mod != nil {
		return mod.Current().Mono
	}
	return 0
}
