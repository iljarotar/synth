package config

type config struct {
	SampleRate      float64
	FadeIn, FadeOut float64
	Duration        int
}

var Default = config{
	SampleRate: 44100,
	FadeIn:     1,
	FadeOut:    5,
	Duration:   0,
}

var Config = config{}
