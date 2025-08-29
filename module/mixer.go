package module

import (
	"fmt"

	"github.com/iljarotar/synth/utils"
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
	m.Gain = utils.Limit(m.Gain, gainLimits[0], gainLimits[1])

	for mod, gain := range m.In {
		m.In[mod] = utils.Limit(gain, gainLimits[0], gainLimits[1])
	}

	return nil
}

func (m *Mixer) Step(modules ModulesMap) {
	var left, right, mono float64

	for name, gain := range m.In {
		if mod, ok := modules[name]; ok {
			left += mod.Current().Left * gain
			right += mod.Current().Right * gain
			mono += mod.Current().Mono * gain
		}
	}

	left = utils.Limit(left*m.Gain, outputLimits[0], outputLimits[1])
	right = utils.Limit(right*m.Gain, outputLimits[0], outputLimits[1])
	mono = utils.Limit(mono*m.Gain, outputLimits[0], outputLimits[1])

	avg := (mono + m.current.Mono) / 2
	m.integral += avg / m.sampleRate

	m.current = Output{
		Mono:  mono,
		Left:  left,
		Right: right,
	}
}
