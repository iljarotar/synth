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
	max, min float64
}

type output struct {
	Mono, Left, Right float64
}

var (
	ampLimits      limits = limits{min: 0, max: 1}
	panLimits      limits = limits{min: -1, max: 1}
	phaseLimits    limits = limits{min: -1, max: 1}
	freqLimits     limits = limits{min: 0, max: 20000}
	bpmLimits      limits = limits{min: 0, max: 600000}
	envelopeLimits limits = limits{min: 0, max: 10000}
)

func modulateValue(modulators []string, modMap ModulesMap) float64 {
	var y float64

	for _, m := range modulators {
		mod, ok := modMap[m]
		if ok {
			y += mod.Current().Mono
		}
	}

	return y
}

func modulate(param Param, lim limits, modMap ModulesMap) float64 {
	y := param.Val + modulateValue(param.Mod, modMap)*param.ModAmp
	return utils.Limit(y, lim.min, lim.max)
}

func stereo(x, pan float64) output {
	out := output{}
	p := utils.Percentage(pan, -1, 1)
	out.Mono = x
	out.Right = x * p
	out.Left = x * (1 - p)

	return out
}
