package wavetable

import (
	"github.com/iljarotar/synth/config"
)

type OscillatorType string

func (t OscillatorType) String() string {
	return string(t)
}

const (
	Sine   OscillatorType = "Sine"
	Square OscillatorType = "Square"
	Noise  OscillatorType = "Noise"
)

type WaveTable struct {
	Step        float64
	Phase       float64 `yaml:"phase"`
	SignalFunc  SignalFunc
	Oscillators []Oscillator `yaml:"oscillators"`
	Filters     []Filter     `yaml:"filters"`
}

type Oscillator struct {
	Type      OscillatorType `yaml:"type"`
	Amplitude float64        `yaml:"amplitude"`
	Freq      float64        `yaml:"freq"`
	PM        *WaveTable     `yaml:"pm"`
	AM        *WaveTable     `yaml:"am"`
}

func (w *WaveTable) Initialize() {
	c := config.Instance()
	w.Step = 1 / c.SampleRate

	f := make([]SignalFunc, 0)

	for i := range w.Filters {
		w.Filters[i].Initialize()
	}

	for i := range w.Oscillators {
		osc := w.Oscillators[i]
		f = append(f, NewSignalFunc(osc.Type))

		if osc.PM != nil {
			osc.PM.Initialize()
		}

		if osc.AM != nil {
			osc.AM.Initialize()
		}
	}

	signalFunc := func(x float64) float64 {
		var y float64

		for i := range w.Oscillators {
			osc := w.Oscillators[i]
			amp := osc.Amplitude
			freq := osc.Freq

			if osc.PM != nil {
				x += osc.PM.SignalFunc(x)
			}

			if osc.AM != nil {
				amp += osc.AM.SignalFunc(x)
			}

			for j := range w.Filters {
				filter := w.Filters[j]
				amp *= filter.Apply(freq, x)
			}

			y += f[i](x*freq) * amp
		}

		return y / float64(len(w.Oscillators))
	}

	w.SignalFunc = signalFunc
}
