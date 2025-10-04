package module

import (
	"math"

	"github.com/iljarotar/synth/calc"
)

type (
	Wavetable struct {
		Module
		Freq       float64   `yaml:"freq"`
		CV         string    `yaml:"cv"`
		Mod        string    `yaml:"mod"`
		Signal     []float64 `yaml:"signal"`
		sampleRate float64
		idx        float64
	}

	WavetableMap map[string]*Wavetable
)

func (m WavetableMap) Initialize(sampleRate float64) {
	for _, w := range m {
		if w == nil {
			continue
		}
		w.initialze(sampleRate)
	}
}

func (w *Wavetable) initialze(sampleRate float64) {
	w.sampleRate = sampleRate
	w.Freq = calc.Limit(w.Freq, freqRange)

	var signal []float64
	for _, x := range w.Signal {
		signal = append(signal, calc.Limit(x, outputRange))
	}
	w.Signal = signal
}

func (w *Wavetable) Step(modules ModuleMap) {
	length := len(w.Signal)
	val := w.Signal[int(math.Floor(w.idx))%length]
	w.current = Output{
		Mono:  val,
		Left:  val / 2,
		Right: val / 2,
	}

	freq := w.Freq
	if w.CV != "" {
		freq = cv(freqRange, getMono(modules[w.CV]))
	}

	mod := math.Pow(2, getMono(modules[w.Mod]))
	w.idx += freq * mod * float64(length) / w.sampleRate
}
