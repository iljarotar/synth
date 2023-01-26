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
	Freq      int            `yaml:"freq"`
}

func (w *WaveTable) CreateSignalFunction() {
	var amp float64
	f := make([]SignalFunc, 0)

	for i := range w.Oscillators {
		w := w.Oscillators[i]
		amp += w.Amplitude
		f = append(f, NewFunc(w.Freq, w.Amplitude, w.Type))
	}

	signalFunc := func(x float64) float64 {
		var y float64
		for i := range w.Oscillators {
			y += f[i](x)
		}
		return y / amp // normalize
	}

	c := config.Instance()
	w.Step = 1 / c.SampleRate
	w.SignalFunc = signalFunc
}
