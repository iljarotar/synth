package module

import (
	"math"

	"github.com/iljarotar/synth/config"
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
	Name   string         `yaml:"name"`
	Type   OscillatorType `yaml:"type"`
	Freq   Param          `yaml:"freq"`
	Amp    Param          `yaml:"amp"`
	Phase  float64        `yaml:"phase"`
	Pan    Param          `yaml:"pan"`
	signal SignalFunc
}

func (o *Oscillator) Initialize() {
	o.signal = newSignalFunc(o.Type)
	o.limitParams()

	y := o.signalValue(0, o.Amp.Val, 0)
	o.current = stereo(y, o.Pan.Val)
}

func (o *Oscillator) Next(t float64, modMap ModulesMap) {
	pan := modulate(o.Pan, panLimits, modMap)
	amp := modulate(o.Amp, ampLimits, modMap)
	offset := o.getOffset(modMap)

	y := o.signalValue(t, amp, offset)
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
	y := o.signal(phi) * amp

	avg := (y + o.Current().Mono) / 2
	o.integral += avg / config.Config.SampleRate

	return y
}

func (o *Oscillator) limitParams() {
	o.Amp.ModAmp = utils.Limit(o.Amp.ModAmp, ampLimits.min, ampLimits.max)
	o.Amp.Val = utils.Limit(o.Amp.Val, ampLimits.min, ampLimits.max)

	o.Phase = utils.Limit(o.Phase, phaseLimits.min, phaseLimits.max)

	o.Pan.ModAmp = utils.Limit(o.Pan.ModAmp, panLimits.min, panLimits.max)
	o.Pan.Val = utils.Limit(o.Pan.Val, panLimits.min, panLimits.max)

	o.Freq.ModAmp = utils.Limit(o.Freq.ModAmp, freqLimits.min, freqLimits.max)
	o.Freq.Val = utils.Limit(o.Freq.Val, freqLimits.min, freqLimits.max)
}
