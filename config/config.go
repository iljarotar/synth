package config

const defaultSampleRate = 44100

type config struct {
	sampleRate float64
}

var Instance = config{sampleRate: defaultSampleRate}

func (c *config) SetSampleRate(sampleRate float64) {
	c.sampleRate = sampleRate
}

func (c *config) SampleRate() float64 {
	return c.sampleRate
}
