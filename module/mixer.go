package module

import (
	"fmt"

	"github.com/iljarotar/synth/calc"
)

type Mixer struct {
	Module
	Gain       float64            `yaml:"gain"`
	CV         string             `yaml:"cv"`
	Mod        string             `yaml:"mod"`
	In         map[string]float64 `yaml:"in"`
	sampleRate float64
}

type MixerMap map[string]*Mixer

func (m MixerMap) Initialize(sampleRate float64) error {
	for name, mixer := range m {
		if err := mixer.initialize(sampleRate); err != nil {
			return fmt.Errorf("failed to initialize mixer %s:%w", name, err)
		}
	}
	return nil
}

func (m *Mixer) initialize(sampleRate float64) error {
	m.sampleRate = sampleRate
	m.Gain = calc.Limit(m.Gain, gainLimits)

	for mod, gain := range m.In {
		m.In[mod] = calc.Limit(gain, gainLimits)
	}

	return nil
}

func (m *Mixer) Step(modules ModulesMap) {
	var (
		left, right, mono float64
	)

	for name, gain := range m.In {
		if mod, ok := modules[name]; ok {
			left += mod.Current().Left * gain
			right += mod.Current().Right * gain
			mono += mod.Current().Mono * gain
		}
	}

	gain := modulate(m.Gain, gainLimits, getMono(modules[m.Mod]))
	if m.CV != "" {
		gain = cv(gainLimits, getMono(modules[m.CV]))
	}

	left = calc.Limit(left*gain, outputLimits)
	right = calc.Limit(right*gain, outputLimits)
	mono = calc.Limit(mono*gain, outputLimits)

	avg := (mono + m.current.Mono) / 2
	m.integral += avg / m.sampleRate

	m.current = Output{
		Mono:  mono,
		Left:  left,
		Right: right,
	}
}
