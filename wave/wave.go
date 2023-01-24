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
}

type Wave struct {
	Type      WaveType
	Amplitude float64
	Freq      int
}

func NewWaveTable(waves ...Wave) WaveTable {
	var amp float64
	f := make([]SignalFunc, 0)

	for i := range waves {
		w := waves[i]
		amp += w.Amplitude
		f = append(f, NewFunc(w.Freq, w.Amplitude, w.Type))
	}

	signalFunc := func(x ...float64) float64 {
		var y float64
		for i := range waves {
			y += f[i](x...)
		}
		return y / amp // normalize
	}

	c := config.Instance()
	return WaveTable{Step: 1 / c.SampleRate, SignalFunc: signalFunc}
}
