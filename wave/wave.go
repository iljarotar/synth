package wave

import (
	"github.com/iljarotar/synth/config"
)

type WaveType string

func (t WaveType) String() string {
	return string(t)
}

const (
	Sine  WaveType = "Sine"
	Noise WaveType = "Noise"
)

type WaveTable struct {
	Step, Phase float64
	SignalFunc  SignalFunc
	Waves       []Wave `yaml:"waves"`
}

type Wave struct {
	Type      WaveType `yaml:"type"`
	Amplitude float64  `yaml:"amplitude"`
	Freq      int      `yaml:"freq"`
}

func (w *WaveTable) CreateSignalFunction() {
	var amp float64
	f := make([]SignalFunc, 0)

	for i := range w.Waves {
		w := w.Waves[i]
		amp += w.Amplitude
		f = append(f, NewFunc(w.Freq, w.Amplitude, w.Type))
	}

	signalFunc := func(x ...float64) float64 {
		var y float64
		for i := range w.Waves {
			y += f[i](x...)
		}
		return y / amp // normalize
	}

	c := config.Instance()
	w.Step = 1 / c.SampleRate
	w.SignalFunc = signalFunc
}
