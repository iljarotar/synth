package module

import (
	"fmt"

	"github.com/iljarotar/synth/calc"
)

type (
	Mixer struct {
		Module
		Gain       float64            `yaml:"gain"`
		CV         string             `yaml:"cv"`
		Mod        string             `yaml:"mod"`
		In         map[string]float64 `yaml:"in"`
		sampleRate float64
	}

	MixerMap map[string]*Mixer
)

func (m MixerMap) Initialize(sampleRate float64) error {
	for name, mixer := range m {
		if mixer == nil {
			continue
		}
		if err := mixer.initialize(sampleRate); err != nil {
			return fmt.Errorf("failed to initialize mixer %s: %w", name, err)
		}
	}
	return nil
}

func (m *Mixer) initialize(sampleRate float64) error {
	m.sampleRate = sampleRate
	m.Gain = calc.Limit(m.Gain, gainRange)

	for mod, gain := range m.In {
		m.In[mod] = calc.Limit(gain, gainRange)
	}

	return nil
}

func (m *Mixer) Update(new *Mixer) {
	if new == nil {
		return
	}

	m.Gain = new.Gain
	m.CV = new.CV
	m.Mod = new.Mod
	m.In = new.In
}

func (m *Mixer) Step(modules ModuleMap) {
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

	gain := m.Gain
	if m.CV != "" {
		gain = cv(gainRange, getMono(modules[m.CV]))
	}
	gain = modulate(gain, gainRange, getMono(modules[m.Mod]))

	left = calc.Limit(left*gain, outputRange)
	right = calc.Limit(right*gain, outputRange)
	mono = calc.Limit(mono*gain, outputRange)

	m.current = Output{
		Mono:  mono,
		Left:  left,
		Right: right,
	}
}
