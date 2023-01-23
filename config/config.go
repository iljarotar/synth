package config

// TODO: make Config a singleton class

type Config struct {
	SampleRate float64
}

func NewConfig(sampleRate float64) *Config {
	return &Config{SampleRate: sampleRate}
}
