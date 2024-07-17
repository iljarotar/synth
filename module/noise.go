package module

import (
	"math/rand"

	"github.com/iljarotar/synth/utils"
)

type Noise struct {
	Module
	Name string `yaml:"name"`
	Amp  Param  `yaml:"amp"`
	Pan  Param  `yaml:"pan"`
}

func (n *Noise) Initialize() {
	n.limitParams()
	n.current = stereo(noise()*n.Amp.Val, n.Pan.Val)
}

func (n *Noise) Next(_ float64, modMap ModulesMap) {
	pan := utils.Limit(n.Pan.Val+modulate(n.Pan.Mod, modMap)*n.Pan.ModAmp, panLimits.min, panLimits.max)
	amp := utils.Limit(n.Amp.Val+modulate(n.Amp.Mod, modMap)*n.Amp.ModAmp, ampLimits.min, ampLimits.max)

	n.current = stereo(noise()*amp, pan)
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
