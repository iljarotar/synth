package config

type Config struct {
	SampleRate      float64
	FadeIn, FadeOut float64
}

var Instance = Config{}
