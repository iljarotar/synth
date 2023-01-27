package oscillator

import (
	"github.com/iljarotar/synth/config"
)

type OscillatorType string

func (t OscillatorType) String() string {
	return string(t)
}

const (
	Sine  OscillatorType = "Sine"
	Noise OscillatorType = "Noise"
)

type WaveTable struct {
	Step, Phase float64
	SignalFunc  SignalFunc
	Oscillators []Oscillator `yaml:"oscillators"`
}

type Oscillator struct {
	Type      OscillatorType `yaml:"type"`
	Amplitude float64        `yaml:"amplitude"`
	Freq      float64        `yaml:"freq"`
	FreqMod   *WaveTable     `yaml:"freq-mod"`
	AmpMod    *WaveTable     `yaml:"amp-mod"`
}

func (w *WaveTable) Initialize() {
	f := make([]SignalFunc, 0)

	for i := range w.Oscillators {
		w := w.Oscillators[i]
		w.Amplitude /= 100 // amplitude should be given in %
		f = append(f, NewFunc(w.Type))

		if w.FreqMod != nil {
			w.FreqMod.Initialize()
		}

		if w.AmpMod != nil {
			w.AmpMod.Initialize()
		}
	}

	signalFunc := func(x float64) float64 {
		var y float64

		for i := range w.Oscillators {
			osc := w.Oscillators[i]
			amp := osc.Amplitude
			freq := osc.Freq

			if osc.FreqMod != nil {
				freq += osc.FreqMod.SignalFunc(x)
			}

			if osc.AmpMod != nil {
				amp += osc.AmpMod.SignalFunc(x)
			}

			y += f[i](x*freq) * amp
		}

		return y / float64(len(w.Oscillators))
	}

	c := config.Instance()
	w.Step = 1 / c.SampleRate
	w.SignalFunc = signalFunc
}
