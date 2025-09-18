package module

import (
	"math"

	"github.com/iljarotar/synth/calc"
)

type Wavetable struct {
	Module
	Freq       float64   `yaml:"freq"`
	CV         string    `yaml:"cv"`
	Mod        string    `yaml:"mod"`
	Signal     []float64 `yaml:"signal"`
	sampleRate float64
	idx        float64
}

type WavetableMap map[string]*Wavetable

func (m WavetableMap) Initialize(sampleRate float64) {
	for _, w := range m {
		w.initialze(sampleRate)
	}
}

func (w *Wavetable) initialze(sampleRate float64) {
	w.sampleRate = sampleRate
	w.Freq = calc.Limit(w.Freq, freqLimits)

	var signal []float64
	for _, x := range w.Signal {
		signal = append(signal, calc.Limit(x, outputLimits))
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
		cv := getMono(modules[w.CV])
		freq = calc.Transpose(cv, outputLimits, freqLimits)
	}

	mod := math.Pow(2, getMono(modules[w.Mod]))
	w.idx += freq * mod * float64(length) / w.sampleRate
}
