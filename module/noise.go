package module

import (
	"math"
	"math/rand"
	"time"

	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/utils"
)

type NoiseMap map[string]*Noise

type Noise struct {
	Name    string  `yaml:"name"`
	Amp     Param   `yaml:"amp"`
	Pan     Param   `yaml:"pan"`
	Low     float64 `yaml:"low"`
	High    float64 `yaml:"high"`
	Current output
	x, b    []float64
}

func (n *Noise) Initialize() {
	rand.Seed(time.Now().Unix())
	n.limitParams()
	n.Current = stereo(noise()*n.Amp.Val, n.Pan.Val)
	n.x = make([]float64, 100)
	n.getFilterCoeffs(n.Low, n.High)
}

func (n *Noise) Next(oscMap OscillatorsMap, customMap CustomMap) {
	pan := utils.Limit(n.Pan.Val+modulate(n.Pan.Mod, oscMap, customMap)*n.Pan.ModAmp, panLimits.low, panLimits.high)
	amp := utils.Limit(n.Amp.Val+modulate(n.Amp.Mod, oscMap, customMap)*n.Amp.ModAmp, ampLimits.low, ampLimits.high)
	n.nextInput(noise())

	n.Current = stereo(n.filteredNoise()*amp, pan)
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

func (n *Noise) getFilterCoeffs(low, high float64) {
	M := len(n.x)
	b := make([]float64, M+1)
	ftLow := low / config.Config.SampleRate
	ftHigh := high / config.Config.SampleRate
	pi := math.Pi

	for i := range b {
		k := float64(i - M/2)
		hamming := 0.54 - 0.46*math.Cos(2*pi*float64(i)/float64(M))

		if k == 0 {
			b[i] = 2 * (ftHigh - ftLow)
		} else {
			b[i] = (math.Sin(2*pi*ftHigh*k) - math.Sin(2*pi*ftLow*k)) / (pi * k)
		}

		b[i] *= hamming
	}

	b = utils.Normalize(b, -0.22, 1) // approximately sinc range

	n.b = b
}

func (n *Noise) filteredNoise() float64 {
	if len(n.b) == 0 {
		return n.x[0]
	}

	y := n.b[0]

	for k := range n.x {
		y += n.b[k+1] * n.x[k]
	}

	return y
}

func (n *Noise) nextInput(x float64) {
	n.x = append([]float64{x}, n.x[:len(n.x)-1]...)
}
