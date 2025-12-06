package module

import (
	"fmt"
	"math"

	"github.com/iljarotar/synth/calc"
)

type (
	Oscillator struct {
		Module
		Type       oscillatorType `yaml:"type"`
		Freq       float64        `yaml:"freq"`
		CV         string         `yaml:"cv"`
		Mod        string         `yaml:"mod"`
		Phase      float64        `yaml:"phase"`
		signal     SignalFunc
		sampleRate float64
		arg        float64
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
	o.Freq = new.Freq
	o.CV = new.CV
	o.Mod = new.Mod
	o.Phase = new.Phase
	o.signal = new.signal
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
}
