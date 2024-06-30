package module

import (
	"math/rand"

	"github.com/iljarotar/synth/utils"
)

type Noise struct {
	Module
	Name   string `yaml:"name"`
	Amp    Param  `yaml:"amp"`
	Pan    Param  `yaml:"pan"`
	Filter Filter `yaml:"filter"`
}

func (n *Noise) Initialize() {
	n.limitParams()
	n.current = stereo(noise()*n.Amp.Val, n.Pan.Val)
	n.Filter.Initialize()
}

func (n *Noise) Next(_ float64, modMap ModulesMap) {
	pan := utils.Limit(n.Pan.Val+modulate(n.Pan.Mod, modMap)*n.Pan.ModAmp, panLimits.min, panLimits.max)
	amp := utils.Limit(n.Amp.Val+modulate(n.Amp.Mod, modMap)*n.Amp.ModAmp, ampLimits.min, ampLimits.max)

	n.current = stereo(n.Filter.Tap(noise())*amp, pan)
}

func (n *Noise) limitParams() {
	n.Amp.ModAmp = utils.Limit(n.Amp.ModAmp, modLimits.min, modLimits.max)
	n.Amp.Val = utils.Limit(n.Amp.Val, ampLimits.min, ampLimits.max)

	n.Pan.ModAmp = utils.Limit(n.Pan.ModAmp, modLimits.min, modLimits.max)
	n.Pan.Val = utils.Limit(n.Pan.Val, panLimits.min, panLimits.max)

	n.Filter.LowCutoff = utils.Limit(n.Filter.LowCutoff, freqLimits.min, freqLimits.max)
	n.Filter.HighCutoff = utils.Limit(n.Filter.HighCutoff, freqLimits.min, freqLimits.max)
}

func noise() float64 {
	return rand.Float64()*2 - 1
}
