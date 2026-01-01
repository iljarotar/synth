package module

import (
	"fmt"

	"github.com/iljarotar/synth/calc"
)

type (
	Mixer struct {
		Module
		Gain float64            `yaml:"gain"`
		CV   string             `yaml:"cv"`
		Mod  string             `yaml:"mod"`
		In   map[string]float64 `yaml:"in"`
		Fade float64            `yaml:"fade"`

		sampleRate float64

		gainFader   *fader
		inputFaders map[string]*fader
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
	m.Fade = calc.Limit(m.Fade, fadeRange)

	m.gainFader = &fader{
		current: m.Gain,
		target:  m.Gain,
	}

	m.inputFaders = map[string]*fader{}
	for mod, gain := range m.In {
		m.In[mod] = calc.Limit(gain, gainRange)

		m.inputFaders[mod] = &fader{
			current: gain,
			target:  gain,
		}
	}
	m.initializeFaders()

	return nil
}

func (m *Mixer) Update(new *Mixer) {
	if new == nil {
		return
	}

	m.CV = new.CV
	m.Mod = new.Mod
	m.Fade = new.Fade
	m.updateGains(new)
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

	if m.gainFader != nil {
		m.Gain = m.gainFader.fade()
	}

	for mod, f := range m.inputFaders {
		if f == nil {
			continue
		}
		m.In[mod] = f.fade()

		if f.current == 0 && f.target == 0 {
			delete(m.inputFaders, mod)
			delete(m.In, mod)
		}
	}
}

func (m *Mixer) initializeFaders() {
	if m.gainFader != nil {
		m.gainFader.initialize(m.Fade, m.sampleRate)
	}

	for _, fader := range m.inputFaders {
		if fader != nil {
			fader.initialize(m.Fade, m.sampleRate)
		}
	}
}

func (m *Mixer) updateGains(new *Mixer) {
	if m.gainFader != nil {
		m.gainFader.target = new.Gain
	}

	for mod, gain := range new.In {
		f, ok := m.inputFaders[mod]

		if ok && f != nil {
			f.target = gain
			continue
		}

		m.inputFaders[mod] = &fader{
			current: 0,
			target:  gain,
		}
		m.In[mod] = 0
	}

	for mod, f := range m.inputFaders {
		if f == nil {
			continue
		}

		if _, ok := new.In[mod]; !ok {
			f.target = 0
		}
	}

	m.initializeFaders()
}
