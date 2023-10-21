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
	pan := utils.Limit(c.Pan.Val+modulate(c.Pan.Mod, modMap)*c.Pan.ModAmp, panLimits.low, panLimits.high)
	amp := utils.Limit(c.Amp.Val+modulate(c.Amp.Mod, modMap)*c.Amp.ModAmp, ampLimits.low, ampLimits.high)
	freq := utils.Limit(c.Freq.Val+modulate(c.Freq.Mod, modMap)*c.Freq.ModAmp, freqLimits.low, freqLimits.high)

	y := c.signalValue(t, amp, freq)
	c.current = stereo(y, pan)
}

func (c *CustomSignal) limitParams() {
	c.Amp.ModAmp = utils.Limit(c.Amp.ModAmp, modLimits.low, modLimits.high)
	c.Amp.Val = utils.Limit(c.Amp.Val, ampLimits.low, ampLimits.high)

	c.Pan.ModAmp = utils.Limit(c.Pan.ModAmp, modLimits.low, modLimits.high)
	c.Pan.Val = utils.Limit(c.Pan.Val, panLimits.low, panLimits.high)

	c.Freq.ModAmp = utils.Limit(c.Freq.ModAmp, freqLimits.low, freqLimits.high)
	c.Freq.Val = utils.Limit(c.Freq.Val, freqLimits.low, freqLimits.high)
}

func (c *CustomSignal) signalValue(t, amp, freq float64) float64 {
	idx := int(math.Floor(t * float64(len(c.Data)) * freq))
	var val float64

	if len(c.Data) > 0 {
		val = c.Data[idx%len(c.Data)]
	}

	y := amp * val
	c.integral += y / config.Config.SampleRate

	return y
}
