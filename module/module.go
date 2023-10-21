package module

import "github.com/iljarotar/synth/utils"

type Module struct {
	integral float64
	current  output
}

type IModule interface {
	Initialize()
	Next(t float64, modMap ModulesMap)
	Integral() float64
	Current() output
}

type ModulesMap map[string]IModule

func (m *Module) Integral() float64 {
	return m.integral
}

func (m *Module) Current() output {
	return m.current
}

type Param struct {
	Val    float64  `yaml:"val"`
	Mod    []string `yaml:"mod"`
	ModAmp float64  `yaml:"mod-amp"`
}

type limits struct {
	high, low float64
}

type output struct {
	Mono, Left, Right float64
}

var (
	ampLimits   limits = limits{low: 0, high: 1}
	modLimits   limits = limits{low: 0, high: 1}
	panLimits   limits = limits{low: -1, high: 1}
	phaseLimits limits = limits{low: -1, high: 1}
	freqLimits  limits = limits{low: 0, high: 20000}
)

func modulate(modulators []string, modMap ModulesMap) float64 {
	var y float64

	for _, m := range modulators {
		mod, ok := modMap[m]
		if ok {
			y += mod.Current().Mono
		}
	}

	return y
}

func stereo(x, pan float64) output {
	out := output{}
	p := utils.Percentage(pan, -1, 1)
	out.Mono = x
	out.Right = x * p
	out.Left = x * (1 - p)

	return out
}
