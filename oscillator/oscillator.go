package oscillator

import (
	"github.com/iljarotar/synth/config"
	f "github.com/iljarotar/synth/filter"
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
	Filters     []f.Filter   `yaml:"filters"`
}

type Oscillator struct {
	Type      OscillatorType `yaml:"type"`
	Amplitude *float64       `yaml:"amplitude"`
	Freq      float64        `yaml:"freq"`
	FreqMod   *WaveTable     `yaml:"freq-mod"`
	AmpMod    *WaveTable     `yaml:"amp-mod"`
}

func (w *WaveTable) Initialize() {
	f := make([]SignalFunc, 0)

	for i := range w.Filters {
		w.Filters[i].Initialize()
	}

	for i := range w.Oscillators {
		w := w.Oscillators[i]
		f = append(f, NewFunc(w.Type))

		if w.FreqMod != nil {
			w.FreqMod.Initialize()
		}

		if w.AmpMod != nil {
			w.AmpMod.Initialize()
		}

		*w.Amplitude /= 100 // amplitude is given in percent
	}

	signalFunc := func(x float64) float64 {
		var y float64

		for i := range w.Oscillators {
			osc := w.Oscillators[i]
			amp := *osc.Amplitude
			freq := osc.Freq

			if osc.FreqMod != nil {
				freq += osc.FreqMod.SignalFunc(x)
			}

			if osc.AmpMod != nil {
				amp += osc.AmpMod.SignalFunc(x)
			}

			// doesn't work
			// for i := range w.Filters {
			// 	filter := w.Filters[i]
			// 	amp += filter.Apply(freq)
			// }

			y += f[i](x*freq) * amp
		}

		return y / float64(len(w.Oscillators))
	}

	c := config.Instance()
	w.Step = 1 / c.SampleRate
	w.SignalFunc = signalFunc
}
