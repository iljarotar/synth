package module

import (
	"math/rand"

	"github.com/iljarotar/synth/utils"
)

type Noise struct {
	Module
	Name    string   `yaml:"name"`
	Amp     Param    `yaml:"amp"`
	Pan     Param    `yaml:"pan"`
	Filters []string `yaml:"filters"`
	inputs  []filterInputs
}

func (n *Noise) Initialize() {
	n.limitParams()
	n.current = stereo(noise()*n.Amp.Val, n.Pan.Val)
	n.inputs = make([]filterInputs, len(n.Filters))
}

func (n *Noise) Next(modMap ModulesMap, filtersMap FiltersMap) {
	pan := modulate(n.Pan, panLimits, modMap)
	amp := modulate(n.Amp, ampLimits, modMap)

	var y2, y float64
	x := noise()

	for i, f := range n.Filters {
		filter, ok := filtersMap[f]
		if !ok {
			continue
		}
		if len(n.inputs) != len(n.Filters) {
			return
		}

		inputs := n.inputs[i]
		y2 = filter.Tap(inputs.x2, inputs.x1, inputs.x0, inputs.y1, inputs.y0)
		y += y2

		inputs.x0 = inputs.x1
		inputs.x1 = inputs.x2
		inputs.y0 = inputs.y1
		inputs.y1 = y2
		inputs.x2 = x
		n.inputs[i] = inputs
	}

	if len(n.Filters) == 0 {
		y = noise()
	} else {
		y /= float64(len(n.Filters))
	}
	n.current = stereo(y*amp, pan)
}

func (n *Noise) limitParams() {
	n.Amp.ModAmp = utils.Limit(n.Amp.ModAmp, ampLimits.min, ampLimits.max)
	n.Amp.Val = utils.Limit(n.Amp.Val, ampLimits.min, ampLimits.max)

	n.Pan.ModAmp = utils.Limit(n.Pan.ModAmp, panLimits.min, panLimits.max)
	n.Pan.Val = utils.Limit(n.Pan.Val, panLimits.min, panLimits.max)
}

func noise() float64 {
	return rand.Float64()*2 - 1
}
