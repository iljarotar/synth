package config

const sampleRate = 44100

var instance *Config = nil

type Config struct {
	SampleRate float64
}

func Instance() *Config {
	if instance == nil {
		instance = &Config{SampleRate: sampleRate}
	}
	return instance
}
