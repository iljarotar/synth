package config

const sampleRate = 44100

var instance *Config = nil

type Config struct {
	SampleRate float64
	RootPath   *string
}

func Instance() *Config {
	rootPath := "examples"
	if instance == nil {
		instance = &Config{SampleRate: sampleRate, RootPath: &rootPath}
	}
	return instance
}
