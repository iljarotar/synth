package module

import (
	"math"

	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/utils"
)

type CustomSignal struct {
	Module
	Name string    `yaml:"name"`
	Data []float64 `yaml:"data"`
	Freq Param     `yaml:"freq"`
	Amp  Param     `yaml:"amp"`
	Pan  Param     `yaml:"pan"`
}

func (c *CustomSignal) Initialize() {
	c.limitParams()
	c.Data = utils.Normalize(c.Data, -1, 1)

	y := c.signalValue(0, c.Amp.Val, c.Freq.Val)
	c.current = stereo(y, c.Pan.Val)
}

func (c *CustomSignal) Next(t float64, modMap ModulesMap) {
	pan := utils.Limit(c.Pan.Val+modulate(c.Pan.Mod, modMap)*c.Pan.ModAmp, panLimits.min, panLimits.max)
	amp := utils.Limit(c.Amp.Val+modulate(c.Amp.Mod, modMap)*c.Amp.ModAmp, ampLimits.min, ampLimits.max)
	freq := utils.Limit(c.Freq.Val+modulate(c.Freq.Mod, modMap)*c.Freq.ModAmp, freqLimits.min, freqLimits.max)

	y := c.signalValue(t, amp, freq)
	c.current = stereo(y, pan)
}

func (c *CustomSignal) limitParams() {
	c.Amp.ModAmp = utils.Limit(c.Amp.ModAmp, modLimits.min, modLimits.max)
	c.Amp.Val = utils.Limit(c.Amp.Val, ampLimits.min, ampLimits.max)

	c.Pan.ModAmp = utils.Limit(c.Pan.ModAmp, modLimits.min, modLimits.max)
	c.Pan.Val = utils.Limit(c.Pan.Val, panLimits.min, panLimits.max)

	c.Freq.ModAmp = utils.Limit(c.Freq.ModAmp, freqLimits.min, freqLimits.max)
	c.Freq.Val = utils.Limit(c.Freq.Val, freqLimits.min, freqLimits.max)
}

func (c *CustomSignal) signalValue(t, amp, freq float64) float64 {
	idx := int(math.Floor(t * float64(len(c.Data)) * freq))
	var val float64

	if l := len(c.Data); l > 0 {
		val = c.Data[idx%l]
	}

	y := amp * val
	c.integral += y / config.Config.SampleRate

	return y
}
