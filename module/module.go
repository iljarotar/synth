package module

import (
	"github.com/iljarotar/synth/calc"
	"github.com/iljarotar/synth/concurrency"
)

type (
	IModule interface {
		Current() Output
	}

	ModuleMap = concurrency.SyncMap[string, IModule]

	Module struct {
		current Output
	}

	Output struct {
		Mono, Left, Right float64
	}
)

var (
	bpmRange = calc.Range{
		Min: 0,
		Max: 2000,
	}
	envelopeRange = calc.Range{
		Min: 1e-15,
		Max: 3600,
	}
	fadeRange = calc.Range{
		Min: 0,
		Max: 3600,
	}
	freqRange = calc.Range{
		Min: 0,
		Max: 20000,
	}
	gainRange = calc.Range{
		Min: 0,
		Max: 1,
	}
	inputGainRange = calc.Range{
		Min: 0,
		Max: 1000,
	}
	outputRange = calc.Range{
		Min: -1,
		Max: 1,
	}
	panRange = calc.Range{
		Min: -1,
		Max: 1,
	}
	pitchRange = calc.Range{
		Min: 400,
		Max: 500,
	}
	transposeRange = calc.Range{
		Min: -24,
		Max: 24,
	}
)

func NewModuleMap(m map[string]IModule) *ModuleMap {
	return concurrency.NewSyncMap(m)
}

func (m *Module) Current() Output {
	return m.current
}

func modulate(x float64, rng calc.Range, mod float64) float64 {
	if mod == 0 {
		return x
	}

	transposed := calc.Transpose(x, rng, outputRange)
	transposed += mod
	transposed = calc.Limit(transposed, outputRange)
	return calc.Transpose(transposed, outputRange, rng)
}

func cv(rng calc.Range, val float64) float64 {
	val = calc.Limit(val, outputRange)
	return calc.Transpose(val, outputRange, rng)
}

func getMono(modules *ModuleMap, name string) float64 {
	mod, found := modules.Get(name)
	if !found || mod == nil {
		return 0
	}
	return mod.Current().Mono
}
