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
	pan     float64
}

func (n *Noise) Initialize() {
	rand.Seed(time.Now().Unix())
	n.limitParams()
	n.Current = stereo(noise(), n.pan)
}

func (n *Noise) Next(oscMap OscillatorsMap) {
	n.pan = utils.Limit(n.Pan.Val+modulate(n.Pan.Mod, oscMap)*n.Pan.ModAmp, panLimits.low, panLimits.high)
	amp := utils.Limit(n.Amp.Val+modulate(n.Amp.Mod, oscMap)*n.Amp.ModAmp, ampLimits.low, ampLimits.high)
	n.Current = stereo(noise()*amp, n.pan)
}

func (n *Noise) limitParams() {
	n.Amp.ModAmp = utils.Limit(n.Amp.ModAmp, ampModLimits.low, ampModLimits.high)
	n.Amp.Val = utils.Limit(n.Amp.Val, ampLimits.low, ampLimits.high)

	n.Pan.ModAmp = utils.Limit(n.Pan.ModAmp, panModLimits.low, panModLimits.high)
	n.Pan.Val = utils.Limit(n.Pan.Val, panLimits.low, panLimits.high)
	n.pan = n.Pan.Val
}

func noise() float64 {
	y := rand.Float64()*2 - 1
	return y
}