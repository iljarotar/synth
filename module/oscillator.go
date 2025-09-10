package module

import (
	"fmt"
	"math"

	"github.com/iljarotar/synth/calc"
)

type OscillatorType string

const (
	Sawtooth        OscillatorType = "Sawtooth"
	ReverseSawtooth OscillatorType = "ReverseSawtooth"
	Sine            OscillatorType = "Sine"
	Square          OscillatorType = "Square"
	Triangle        OscillatorType = "Triangle"
)

type Oscillator struct {
	Module
	Type       OscillatorType `yaml:"type"`
	Freq       float64        `yaml:"freq"`
	CV         string         `yaml:"cv"`
	Mod        string         `yaml:"mod"`
	Phase      float64        `yaml:"phase"`
	signal     SignalFunc
	sampleRate float64
	arg        float64
}

type OscillatorMap map[string]*Oscillator

func (m OscillatorMap) Initialize(sampleRate float64) error {
	for name, osc := range m {
		if err := osc.initialize(sampleRate); err != nil {
			return fmt.Errorf("failed to initialze oscillator %s:%w", name, err)
		}
	}
	return nil
}

func (o *Oscillator) initialize(sampleRate float64) error {
	o.sampleRate = sampleRate
	o.Freq = calc.Limit(o.Freq, freqLimits)

	signal, err := newSignalFunc(o.Type)
	if err != nil {
		return err
	}
	o.signal = signal

	return nil
}

func (o *Oscillator) Step(t float64, modules ModulesMap) {
	freq := o.Freq
	if o.CV != "" {
		cv := getMono(modules[o.CV])
		freq = calc.Transpose(cv, outputLimits, freqLimits)
	}

	c := 2 * math.Pi * o.Phase
	mod := math.Pow(2, getMono(modules[o.Mod]))
	o.arg += 2 * math.Pi * freq * mod / o.sampleRate

	val := o.signal(o.arg + c)
	o.current = Output{
		Mono:  val,
		Left:  val / 2,
		Right: val / 2,
	}
}
