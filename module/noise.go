package module

import (
	"math/rand"

	"github.com/iljarotar/synth/utils"
)

type Noise struct {
	Module
	Name       string   `yaml:"name"`
	Amp        Input    `yaml:"amp"`
	Pan        Input    `yaml:"pan"`
	Filters    []string `yaml:"filters"`
	Envelope   string   `yaml:"envelope"`
	inputs     []filterInputs
	sampleRate float64
}

func (n *Noise) Initialize(sampleRate float64) {
	n.sampleRate = sampleRate
	n.limitParams()
	n.inputs = make([]filterInputs, len(n.Filters))
	n.current = stereo(noise()*n.Amp.Val, n.Pan.Val)
}

func (n *Noise) Next(modMap ModulesMap, filtersMap FiltersMap, envelopesMap EnvelopesMap) {
	pan := modulate(n.Pan, panLimits, modMap)
	amp := modulate(n.Amp, ampLimits, modMap)

	cfg := filterConfig{
		filterNames: n.Filters,
		inputs:      n.inputs,
		FiltersMap:  filtersMap,
	}

	y, newInputs := cfg.applyFilters(noise())
	y = applyEnvelope(y, n.Envelope, envelopesMap)
	n.integral += y / n.sampleRate
	n.inputs = newInputs
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
