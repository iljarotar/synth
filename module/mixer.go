package module

import (
	"fmt"

	"github.com/iljarotar/synth/calc"
	"github.com/iljarotar/synth/concurrency"
)

type (
	Mixer struct {
		Module
		Gain float64            `yaml:"gain"`
		CV   string             `yaml:"cv"`
		Mod  string             `yaml:"mod"`
		In   map[string]float64 `yaml:"in"`
		Fade float64            `yaml:"fade"`

		in         *concurrency.SyncMap[string, float64]
		sampleRate float64

		gainFader   *fader
		inputFaders *concurrency.SyncMap[string, *fader]
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

	m.inputFaders = concurrency.NewSyncMap(map[string]*fader{})
	for mod, gain := range m.In {
		m.In[mod] = calc.Limit(gain, inputGainRange)

		m.inputFaders.Set(mod, &fader{
			current: gain,
			target:  gain,
		})
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

func (m *Mixer) Step(modules *ModuleMap) {
	var (
		left, right, mono float64
	)

	for name, gain := range m.In {
		if mod := modules.Get(name); mod != nil {
			left += mod.Current().Left * gain
			right += mod.Current().Right * gain
			mono += mod.Current().Mono * gain
		}
	}

	gain := m.Gain
	if m.CV != "" {
		gain = cv(gainRange, getMono(modules, m.CV))
	}
	gain = modulate(gain, gainRange, getMono(modules, m.Mod))

	left = calc.Limit(left*gain, outputRange)
	right = calc.Limit(right*gain, outputRange)
	mono = calc.Limit(mono*gain, outputRange)

	m.current = Output{
		Mono:  mono,
		Left:  left,
		Right: right,
	}

	m.fade()
}

func (m *Mixer) fade() {
	if m.gainFader != nil {
		m.Gain = m.gainFader.fade()
	}

	for _, name := range m.inputFaders.Keys() {
		f := m.inputFaders.Get(name)
		if f == nil {
			continue
		}
		m.In[name] = f.fade()

		if f.current == 0 && f.target == 0 {
			m.inputFaders.Delete(name)
			delete(m.In, name)
		}
	}
}

func (m *Mixer) initializeFaders() {
	if m.gainFader != nil {
		m.gainFader.initialize(m.Fade, m.sampleRate)
	}

	for _, name := range m.inputFaders.Keys() {
		if f := m.inputFaders.Get(name); f != nil {
			f.initialize(m.Fade, m.sampleRate)
		}
	}
}

func (m *Mixer) updateGains(new *Mixer) {
	if m.gainFader != nil {
		m.gainFader.target = new.Gain
	}

	if m.inputFaders == nil {
		m.inputFaders = concurrency.NewSyncMap(map[string]*fader{})
	}
	if m.In == nil {
		m.In = map[string]float64{}
	}

	for mod, gain := range new.In {
		f := m.inputFaders.Get(mod)
		if f != nil {
			f.target = gain
			continue
		}

		m.inputFaders.Set(mod, &fader{
			current: 0,
			target:  gain,
		})
		m.In[mod] = 0
	}

	for _, name := range m.inputFaders.Keys() {
		f := m.inputFaders.Get(name)
		if f == nil {
			continue
		}

		if _, ok := new.In[name]; !ok {
			f.target = 0
		}
	}

	m.initializeFaders()
}
