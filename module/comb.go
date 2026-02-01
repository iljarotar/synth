package module

import (
	"math"

	"github.com/iljarotar/synth/calc"
)

type (
	comb struct {
		y          []float64
		sampleRate float64
		idx        int
	}
)

func (c *comb) initialize(seconds, sampleRate float64) {
	c.sampleRate = sampleRate
	length := int(math.Ceil(seconds * sampleRate))
	c.y = make([]float64, length)
}

func (c *comb) update(seconds float64) {
	length := int(math.Ceil(seconds * c.sampleRate))
	if length == len(c.y) {
		return
	}

	if length > len(c.y) {
		diff := length - len(c.y)
		c.y = append(c.y, make([]float64, diff)...)
		return
	}

	c.y = c.y[:length]
	c.idx = int(calc.Limit(float64(c.idx), calc.Range{
		Min: 0,
		Max: float64(len(c.y) - 1),
	}))
}

func (c *comb) step(x, mix float64) float64 {
	y := x * (1 - mix)

	if len(c.y) > 0 {
		y += c.y[c.idx] * mix
		c.y[c.idx] = y
		c.idx = (c.idx + 1) % len(c.y)
	}

	return y
}
