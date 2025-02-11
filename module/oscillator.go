package module

import (
	"math"

	"github.com/iljarotar/synth/utils"
)

type OscillatorType string

func (t OscillatorType) String() string {
	return string(t)
}

const (
	Sawtooth        OscillatorType = "Sawtooth"
	ReverseSawtooth OscillatorType = "ReverseSawtooth"
	Sine            OscillatorType = "Sine"
	Square          OscillatorType = "Square"
	Triangle        OscillatorType = "Triangle"
)

type Oscillator struct {
	Module
	Name       string         `yaml:"name"`
	Type       OscillatorType `yaml:"type"`
	Freq       Input          `yaml:"freq"`
	Amp        Input          `yaml:"amp"`
	Phase      float64        `yaml:"phase"`
	Pan        Input          `yaml:"pan"`
	Filters    []string       `yaml:"filters"`
	Envelope   *Envelope      `yaml:"envelope"`
	inputs     []filterInputs
	signal     SignalFunc
	sampleRate float64
}

func (o *Oscillator) Initialize(sampleRate float64) error {
	if o.Envelope != nil {
		o.Envelope.Initialize()
	}

	o.sampleRate = sampleRate
	signal, err := newSignalFunc(o.Type)
	if err != nil {
		return err
	}
	o.signal = signal

	o.limitParams()
	o.inputs = make([]filterInputs, len(o.Filters))

	y := o.signalValue(0, o.Amp.Val, 0)
	o.current = stereo(y, o.Pan.Val)

	return nil
}

func (o *Oscillator) Next(t float64, modMap ModulesMap, filtersMap FiltersMap) {
	if o.Envelope != nil {
		o.Envelope.Next(t, modMap)
	}

	pan := modulate(o.Pan, panLimits, modMap)
	amp := modulate(o.Amp, ampLimits, modMap)
	offset := o.getOffset(modMap)

	cfg := filterConfig{
		filterNames: o.Filters,
		inputs:      o.inputs,
		FiltersMap:  filtersMap,
	}

	x := o.signalValue(t, amp, offset)
	y, newInputs := cfg.applyFilters(x)
	y = applyEnvelope(y, o.Envelope)
	avg := (y + o.Current().Mono) / 2
	o.integral += avg / o.sampleRate
	o.inputs = newInputs
	o.current = stereo(y, pan)
}

func (o *Oscillator) getOffset(modMap ModulesMap) float64 {
	var y float64

	for _, m := range o.Freq.Mod {
		mod, ok := modMap[m]
		if ok {
			y += mod.Integral()
		}
	}

	return y * o.Freq.ModAmp
}

func (o *Oscillator) signalValue(t, amp, offset float64) float64 {
	shift := o.Phase / o.Freq.Val // shift is a fraction of one period
	phi := 2 * math.Pi * (o.Freq.Val*(t+shift) + offset)
	return o.signal(phi) * amp
}

func (o *Oscillator) limitParams() {
	o.Amp.ModAmp = utils.Limit(o.Amp.ModAmp, ampLimits.min, ampLimits.max)
	o.Amp.Val = utils.Limit(o.Amp.Val, -ampLimits.max, ampLimits.max)

	o.Phase = utils.Limit(o.Phase, phaseLimits.min, phaseLimits.max)

	o.Pan.ModAmp = utils.Limit(o.Pan.ModAmp, panLimits.min, panLimits.max)
	o.Pan.Val = utils.Limit(o.Pan.Val, panLimits.min, panLimits.max)

	o.Freq.ModAmp = utils.Limit(o.Freq.ModAmp, freqLimits.min, freqLimits.max)
	o.Freq.Val = utils.Limit(o.Freq.Val, -freqLimits.max, freqLimits.max)
}
