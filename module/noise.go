package module

import (
	"math/rand"
	"time"

	"github.com/iljarotar/synth/utils"
)

type NoiseMap map[string]*Noise

type Noise struct {
	Name    string `yaml:"name"`
	Amp     Param  `yaml:"amp"`
	Pan     Param  `yaml:"pan"`
	Current output
}

func (n *Noise) Initialize() {
	rand.Seed(time.Now().Unix())
	n.limitParams()
	n.Current = stereo(noise()*n.Amp.Val, n.Pan.Val)
}

func (n *Noise) Next(oscMap OscillatorsMap) {
	pan := utils.Limit(n.Pan.Val+modulate(n.Pan.Mod, oscMap)*n.Pan.ModAmp, panLimits.low, panLimits.high)
	amp := utils.Limit(n.Amp.Val+modulate(n.Amp.Mod, oscMap)*n.Amp.ModAmp, ampLimits.low, ampLimits.high)
	n.Current = stereo(noise()*amp, pan)
}

func (n *Noise) limitParams() {
	n.Amp.ModAmp = utils.Limit(n.Amp.ModAmp, modLimits.low, modLimits.high)
	n.Amp.Val = utils.Limit(n.Amp.Val, ampLimits.low, ampLimits.high)

	n.Pan.ModAmp = utils.Limit(n.Pan.ModAmp, modLimits.low, modLimits.high)
	n.Pan.Val = utils.Limit(n.Pan.Val, panLimits.low, panLimits.high)
}

func noise() float64 {
	return rand.Float64()*2 - 1
}