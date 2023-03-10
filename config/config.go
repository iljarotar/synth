package config

import "math"

const maxSamples = 317520000

type config struct {
	SampleRate      float64
	FadeIn, FadeOut float64
	Duration        float64
	MaxDuration     float64
}

var Default = config{
	SampleRate: 44100,
	FadeIn:     1,
	FadeOut:    1,
	Duration:   -1,
}

var Config = config{}

func (c *config) GetMaxDuration() float64 {
	return math.Floor(maxSamples/c.SampleRate - c.FadeIn - c.FadeOut)
}
