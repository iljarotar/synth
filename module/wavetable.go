package module

import (
	"math"

	"github.com/iljarotar/synth/calc"
)

type (
	Wavetable struct {
		Module
		Freq   float64   `yaml:"freq"`
		CV     string    `yaml:"cv"`
		Mod    string    `yaml:"mod"`
		Signal []float64 `yaml:"signal"`
		Fade   float64   `yaml:"fade"`

		sampleRate float64
		idx        float64

		freqFader *fader
	}

	WavetableMap map[string]*Wavetable
)

func (m WavetableMap) Initialize(sampleRate float64) {
	for _, w := range m {
		if w == nil {
			continue
		}
		w.initialize(sampleRate)
	}
}

func (w *Wavetable) initialize(sampleRate float64) {
	w.sampleRate = sampleRate
	w.Freq = calc.Limit(w.Freq, freqRange)
	w.Fade = calc.Limit(w.Fade, fadeRange)

	w.freqFader = &fader{
		current: w.Freq,
		target:  w.Freq,
	}
	w.freqFader.initialize(w.Fade, sampleRate)

	var signal []float64
	for _, x := range w.Signal {
		signal = append(signal, calc.Limit(x, outputRange))
	}
	w.Signal = signal
}

func (w *Wavetable) Update(new *Wavetable) {
	if new == nil {
		return
	}

	w.CV = new.CV
	w.Mod = new.Mod
	w.Signal = new.Signal
	w.Fade = new.Fade

	if w.freqFader != nil {
		w.freqFader.target = new.Freq
		w.freqFader.initialize(w.Fade, w.sampleRate)
	}
}

func (w *Wavetable) Step(modules *ModuleMap) {
	if len(w.Signal) < 1 {
		return
	}
	val := w.Signal[int(math.Floor(w.idx))%len(w.Signal)]
	w.current = Output{
		Mono:  val,
		Left:  val / 2,
		Right: val / 2,
	}

	freq := w.Freq
	if w.CV != "" {
		freq = cv(freqRange, getMono(modules, w.CV))
	}

	mod := math.Pow(2, getMono(modules, w.Mod))
	w.idx += freq * mod * float64(len(w.Signal)) / w.sampleRate

	w.fade()
}

func (w *Wavetable) fade() {
	if w.freqFader != nil {
		w.Freq = w.freqFader.fade()
	}
}
