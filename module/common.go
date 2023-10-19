package module

import "github.com/iljarotar/synth/utils"

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

func modulate(modulators []string, oscMap OscillatorsMap, customMap CustomMap) float64 {
	var y float64

	for _, mod := range modulators {
		osc, ok := oscMap[mod]
		if ok {
			y += osc.Current.Mono
		}

		c, ok := customMap[mod]
		if ok {
			y += c.Current.Mono
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
