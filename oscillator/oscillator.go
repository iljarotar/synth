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
	Amplitude float64        `yaml:"amplitude"`
	Freq      float64        `yaml:"freq"`
	FM        *WaveTable     `yaml:"fm"`
	AM        *WaveTable     `yaml:"am"`
}

func (w *WaveTable) Initialize() {
	f := make([]SignalFunc, 0)

	for i := range w.Filters {
		w.Filters[i].Initialize()
	}

	for i := range w.Oscillators {
		w := w.Oscillators[i]
		f = append(f, NewFunc(w.Type))

		if w.FM != nil {
			w.FM.Initialize()
		}

		if w.AM != nil {
			w.AM.Initialize()
		}

		w.Amplitude /= 100 // amplitude is given in percent
	}

	signalFunc := func(x float64) float64 {
		var y float64

		for i := range w.Oscillators {
			osc := w.Oscillators[i]
			amp := osc.Amplitude
			freq := osc.Freq

			if osc.FM != nil {
				freq += osc.FM.SignalFunc(x)
			}

			if osc.AM != nil {
				amp += osc.AM.SignalFunc(x)
			}

			for j := range w.Filters {
				filter := w.Filters[j]
				amp *= filter.Apply(osc.Freq)
			}

			y += f[i](x*freq) * amp
		}

		return y / float64(len(w.Oscillators))
	}

	c := config.Instance()
	w.Step = 1 / c.SampleRate
	w.SignalFunc = signalFunc
}
