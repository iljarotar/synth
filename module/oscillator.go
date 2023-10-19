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

type OscillatorsMap map[string]*Oscillator

type Oscillator struct {
	Name     string         `yaml:"name"`
	Type     OscillatorType `yaml:"type"`
	Freq     Param          `yaml:"freq"`
	Amp      Param          `yaml:"amp"`
	Phase    float64        `yaml:"phase"`
	Pan      Param          `yaml:"pan"`
	signal   SignalFunc
	Integral float64
	Current  output
}

func (o *Oscillator) Initialize() {
	o.signal = newSignalFunc(o.Type)
	o.limitParams()

	y := o.signalValue(0, o.Amp.Val, 0)
	o.Current = stereo(y, o.Pan.Val)
}

func (o *Oscillator) Next(x float64, oscMap OscillatorsMap, customMap CustomMap) {
	pan := utils.Limit(o.Pan.Val+modulate(o.Pan.Mod, oscMap, customMap)*o.Pan.ModAmp, panLimits.low, panLimits.high)
	amp := utils.Limit(o.Amp.Val+modulate(o.Amp.Mod, oscMap, customMap)*o.Amp.ModAmp, ampLimits.low, ampLimits.high)
	offset := o.getOffset(oscMap, customMap)

	y := o.signalValue(x, amp, offset)
	o.Current = stereo(y, pan)
}

func (o *Oscillator) getOffset(oscMap OscillatorsMap, customMap CustomMap) float64 {
	var y float64

	for _, mod := range o.Freq.Mod {
		osc, ok := oscMap[mod]
		if ok {
			y += osc.Integral
		}

		c, ok := customMap[mod]
		if ok {
			y += c.Integral
		}
	}

	return y * o.Freq.ModAmp
}

func (o *Oscillator) signalValue(x, amp, offset float64) float64 {
	shift := o.Phase / o.Freq.Val // shift is a fraction of one period
	phi := 2 * math.Pi * (o.Freq.Val*(x+shift) + offset)
	y := o.signal(phi) * amp

	avg := (y + o.Current.Mono) / 2
	o.Integral += avg / config.Config.SampleRate

	return y
}

func (o *Oscillator) limitParams() {
	o.Amp.ModAmp = utils.Limit(o.Amp.ModAmp, modLimits.low, modLimits.high)
	o.Amp.Val = utils.Limit(o.Amp.Val, ampLimits.low, ampLimits.high)

	o.Phase = utils.Limit(o.Phase, phaseLimits.low, phaseLimits.high)

	o.Pan.ModAmp = utils.Limit(o.Pan.ModAmp, modLimits.low, modLimits.high)
	o.Pan.Val = utils.Limit(o.Pan.Val, panLimits.low, panLimits.high)

	o.Freq.ModAmp = utils.Limit(o.Freq.ModAmp, freqLimits.low, freqLimits.high)
	o.Freq.Val = utils.Limit(o.Freq.Val, freqLimits.low, freqLimits.high)
}
