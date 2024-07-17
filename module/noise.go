package module

import (
	"math/rand"

	"github.com/iljarotar/synth/utils"
)

type Noise struct {
	Module
	Name               string   `yaml:"name"`
	Amp                Param    `yaml:"amp"`
	Pan                Param    `yaml:"pan"`
	Filters            []string `yaml:"filters"`
	x2, x1, x0, y1, y0 float64
}

func (n *Noise) Initialize() {
	n.limitParams()
	n.current = stereo(noise()*n.Amp.Val, n.Pan.Val)
}

func (n *Noise) Next(modMap ModulesMap, filtersMap FiltersMap) {
	pan := modulate(n.Pan, panLimits, modMap)
	amp := modulate(n.Amp, ampLimits, modMap)

	var y2 float64

	for _, f := range n.Filters {
		filter, ok := filtersMap[f]
		if !ok {
			continue
		}

		// CONTINUE: x values for cascading taps must be y values of previous ones... maybe recursive?
		y2 := filter.Tap(n.x2, n.x1, n.x0, n.y1, n.y0)

		n.x0 = n.x1
		n.x1 = n.x2
		n.x2 = noise()
		n.y0 = n.y1
		n.y1 = y2
	}

	n.current = stereo(y2*amp, pan)
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
