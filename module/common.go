package module

import "github.com/iljarotar/synth/utils"

type limits struct {
	high, low float64
}

var (
	ampLimits   limits = limits{low: 0, high: 1}
	modLimits   limits = limits{low: 0, high: 1}
	panLimits   limits = limits{low: -1, high: 1}
	phaseLimits limits = limits{low: -1, high: 1}
	freqLimits  limits = limits{low: 0, high: 20000}
)

type output struct {
	Mono, Left, Right float64
}

func modulate(modulators []string, oscMap OscillatorsMap) float64 {
	var y float64

	for _, m := range modulators {
		mod, ok := oscMap[m]
		if ok {
			y += mod.Current.Mono
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
