package module

import (
	"math/rand"
	"time"

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
	rand.Seed(time.Now().Unix())
	n.limitParams()
	n.current = stereo(noise()*n.Amp.Val, n.Pan.Val)
	n.Filter.Initialize()
}

func (n *Noise) Next(_ float64, modMap ModulesMap) {
	pan := utils.Limit(n.Pan.Val+modulate(n.Pan.Mod, modMap)*n.Pan.ModAmp, panLimits.low, panLimits.high)
	amp := utils.Limit(n.Amp.Val+modulate(n.Amp.Mod, modMap)*n.Amp.ModAmp, ampLimits.low, ampLimits.high)

	n.current = stereo(n.Filter.Tap(noise())*amp, pan)
}

func (n *Noise) limitParams() {
	n.Amp.ModAmp = utils.Limit(n.Amp.ModAmp, modLimits.low, modLimits.high)
	n.Amp.Val = utils.Limit(n.Amp.Val, ampLimits.low, ampLimits.high)

	n.Pan.ModAmp = utils.Limit(n.Pan.ModAmp, modLimits.low, modLimits.high)
	n.Pan.Val = utils.Limit(n.Pan.Val, panLimits.low, panLimits.high)

	n.Filter.LowCutoff = utils.Limit(n.Filter.LowCutoff, freqLimits.low, freqLimits.high)
	n.Filter.HighCutoff = utils.Limit(n.Filter.HighCutoff, freqLimits.low, freqLimits.high)
}

func noise() float64 {
	return rand.Float64()*2 - 1
}
