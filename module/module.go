package module

import (
	"github.com/iljarotar/synth/utils"
)

type Module struct {
	integral float64
	current  output
}

type IModule interface {
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

type Input struct {
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
	ampLimits       limits = limits{min: 0, max: 2}
	bpmLimits       limits = limits{min: 0.0001, max: 600000}
	cutoffLimits    limits = limits{min: 1, max: 20000}
	envelopeLimits  limits = limits{min: 0, max: 10000}
	freqLimits      limits = limits{min: 0, max: 20000}
	panLimits       limits = limits{min: -1, max: 1}
	phaseLimits     limits = limits{min: -1, max: 1}
	pitchLimits     limits = limits{min: 400, max: 500}
	delayLimits     limits = limits{min: 0, max: 3600}
	transposeLimits limits = limits{min: -24, max: 24}
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

func modulate(param Input, lim limits, modMap ModulesMap) float64 {
	y := param.Val + modulateValue(param.Mod, modMap)*param.ModAmp
	return utils.Limit(y, lim.min, lim.max)
}

func applyEnvelope(x float64, envelope *Envelope) float64 {
	if envelope != nil {
		return x * envelope.current
	}
	return x
}

func stereo(x, pan float64) output {
	out := output{}
	p := utils.Percentage(pan, -1, 1)
	out.Mono = x
	out.Right = x * p
	out.Left = x * (1 - p)

	return out
}
