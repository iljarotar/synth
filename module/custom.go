package module

import (
	"math"

	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/utils"
)

type CustomMap map[string]*Custom

type Custom struct {
	Name     string    `yaml:"name"`
	Data     []float64 `yaml:"data"`
	Freq     Param     `yaml:"freq"`
	Amp      Param     `yaml:"amp"`
	Pan      Param     `yaml:"pan"`
	Integral float64
	Current  output
}

func (c *Custom) Initialize() {
	c.limitParams()
	c.Data = utils.Normalize(c.Data, -1, 1)

	y := c.signalValue(0, c.Amp.Val, c.Freq.Val)
	c.Current = stereo(y, c.Pan.Val)
}

func (c *Custom) Next(x float64, oscMap OscillatorsMap, cMap CustomMap) {
	pan := utils.Limit(c.Pan.Val+modulate(c.Pan.Mod, oscMap, cMap)*c.Pan.ModAmp, panLimits.low, panLimits.high)
	amp := utils.Limit(c.Amp.Val+modulate(c.Amp.Mod, oscMap, cMap)*c.Amp.ModAmp, ampLimits.low, ampLimits.high)
	freq := utils.Limit(c.Freq.Val+modulate(c.Freq.Mod, oscMap, cMap)*c.Freq.ModAmp, freqLimits.low, freqLimits.high)

	y := c.signalValue(x, amp, freq)
	c.Current = stereo(y, pan)
}

func (c *Custom) limitParams() {
	c.Amp.ModAmp = utils.Limit(c.Amp.ModAmp, modLimits.low, modLimits.high)
	c.Amp.Val = utils.Limit(c.Amp.Val, ampLimits.low, ampLimits.high)

	c.Pan.ModAmp = utils.Limit(c.Pan.ModAmp, modLimits.low, modLimits.high)
	c.Pan.Val = utils.Limit(c.Pan.Val, panLimits.low, panLimits.high)

	c.Freq.ModAmp = utils.Limit(c.Freq.ModAmp, freqLimits.low, freqLimits.high)
	c.Freq.Val = utils.Limit(c.Freq.Val, freqLimits.low, freqLimits.high)
}

func (c *Custom) signalValue(x, amp, freq float64) float64 {
	idx := int(math.Floor(x * freq))
	val := c.Data[idx%len(c.Data)]
	y := amp * val

	c.Integral += y / config.Config.SampleRate

	return y
}
