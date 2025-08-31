package module

import "github.com/iljarotar/synth/calc"

type IModule interface {
	Current() Output
	Integral() float64
}

type ModulesMap map[string]IModule

type Module struct {
	current  Output
	integral float64
}

func (m *Module) Current() Output {
	return m.current
}

func (m *Module) Integral() float64 {
	return m.integral
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
)

type Output struct {
	Mono, Left, Right float64
}

func modulate(x float64, rng calc.Range, val float64) float64 {
	transposed := calc.Transpose(x, rng, outputLimits)
	transposed += val
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

func getIntegral(mod IModule) float64 {
	if mod != nil {
		return mod.Integral()
	}
	return 0
}
