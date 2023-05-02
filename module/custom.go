package module

import "github.com/iljarotar/synth/utils"

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
}

func (c *Custom) limitParams() {
	c.Amp.ModAmp = utils.Limit(c.Amp.ModAmp, modLimits.low, modLimits.high)
	c.Amp.Val = utils.Limit(c.Amp.Val, ampLimits.low, ampLimits.high)

	c.Pan.ModAmp = utils.Limit(c.Pan.ModAmp, modLimits.low, modLimits.high)
	c.Pan.Val = utils.Limit(c.Pan.Val, panLimits.low, panLimits.high)

	c.Freq.ModAmp = utils.Limit(c.Freq.ModAmp, freqLimits.low, freqLimits.high)
	c.Freq.Val = utils.Limit(c.Freq.Val, freqLimits.low, freqLimits.high)
}

func (c *Custom) normalize() {
	var min, max float64

	for _, d := range c.Data {
		if d > max {
			max = d
			continue
		}
		if d < min {
			min = d
		}
	}

	// r := max - min

	// for _, d := range c.Data {
	// 	d =
	// }
}
