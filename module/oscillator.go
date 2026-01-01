package module

import (
	"fmt"
	"math"

	"github.com/iljarotar/synth/calc"
)

type (
	Oscillator struct {
		Module
		Type  oscillatorType `yaml:"type"`
		Freq  float64        `yaml:"freq"`
		CV    string         `yaml:"cv"`
		Mod   string         `yaml:"mod"`
		Phase float64        `yaml:"phase"`
		Fade  float64        `yaml:"fade"`

		signal     SignalFunc
		sampleRate float64
		arg        float64

		freqFader  *fader
		phaseFader *fader
	}

	OscillatorMap  map[string]*Oscillator
	oscillatorType string
)

const (
	oscillatorTypeSawtooth        oscillatorType = "Sawtooth"
	oscillatorTypeReverseSawtooth oscillatorType = "ReverseSawtooth"
	oscillatorTypeSine            oscillatorType = "Sine"
	oscillatorTypeSquare          oscillatorType = "Square"
	oscillatorTypeTriangle        oscillatorType = "Triangle"
)

func (m OscillatorMap) Initialize(sampleRate float64) error {
	for name, o := range m {
		if o == nil {
			continue
		}
		if err := o.initialize(sampleRate); err != nil {
			return fmt.Errorf("failed to initialze oscillator %s: %w", name, err)
		}
	}
	return nil
}

func (o *Oscillator) initialize(sampleRate float64) error {
	o.sampleRate = sampleRate
	o.Freq = calc.Limit(o.Freq, freqRange)
	o.Fade = calc.Limit(o.Fade, fadeRange)

	o.freqFader = &fader{
		current: o.Freq,
		target:  o.Freq,
	}
	o.phaseFader = &fader{
		current: o.Phase,
		target:  o.Phase,
	}
	o.initializeFaders()

	signal, err := newSignalFunc(o.Type)
	if err != nil {
		return err
	}
	o.signal = signal

	return nil
}

func (o *Oscillator) Update(new *Oscillator) {
	if new == nil {
		return
	}

	o.Type = new.Type
	o.CV = new.CV
	o.Mod = new.Mod
	o.Fade = new.Fade
	o.signal = new.signal

	if o.freqFader != nil {
		o.freqFader.target = new.Freq
	}
	if o.phaseFader != nil {
		o.phaseFader.target = new.Phase
	}
	o.initializeFaders()
}

func (o *Oscillator) Step(modules ModuleMap) {
	twoPi := 2 * math.Pi
	freq := o.Freq
	if o.CV != "" {
		freq = cv(freqRange, getMono(modules[o.CV]))
	}

	c := twoPi * o.Phase
	mod := math.Pow(2, getMono(modules[o.Mod]))

	val := o.signal(o.arg + c)
	o.current = Output{
		Mono:  val,
		Left:  val / 2,
		Right: val / 2,
	}

	o.arg += twoPi * freq * mod / o.sampleRate
	if o.freqFader != nil {
		o.Freq = o.freqFader.fade()
	}
	if o.phaseFader != nil {
		o.Phase = o.phaseFader.fade()
	}
}

func (o *Oscillator) initializeFaders() {
	if o.freqFader != nil {
		o.freqFader.initialize(o.Fade, o.sampleRate)
	}
	if o.phaseFader != nil {
		o.phaseFader.initialize(o.Fade, o.sampleRate)
	}
}
