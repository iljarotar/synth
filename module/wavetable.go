package module

import (
	"math"

	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/utils"
)

type Wavetable struct {
	Module
	Name    string    `yaml:"name"`
	Table   []float64 `yaml:"table"`
	Freq    Input     `yaml:"freq"`
	Amp     Input     `yaml:"amp"`
	Pan     Input     `yaml:"pan"`
	Filters []string  `yaml:"filters"`
	inputs  []filterInputs
}

func (w *Wavetable) Initialize() {
	w.limitParams()
	w.Table = utils.Normalize(w.Table, -1, 1)
	w.inputs = make([]filterInputs, len(w.Filters))

	y := w.signalValue(0, w.Amp.Val, w.Freq.Val)
	w.current = stereo(y, w.Pan.Val)
}

func (w *Wavetable) Next(t float64, modMap ModulesMap, filtersMap FiltersMap) {
	pan := modulate(w.Pan, panLimits, modMap)
	amp := modulate(w.Amp, ampLimits, modMap)
	freq := modulate(w.Freq, freqLimits, modMap)

	cfg := filterConfig{
		filterNames: w.Filters,
		inputs:      w.inputs,
		FiltersMap:  filtersMap,
	}

	x := w.signalValue(t, amp, freq)
	y, newInputs := cfg.applyFilters(x)
	w.inputs = newInputs
	w.current = stereo(y, pan)
}

func (w *Wavetable) limitParams() {
	w.Amp.ModAmp = utils.Limit(w.Amp.ModAmp, ampLimits.min, ampLimits.max)
	w.Amp.Val = utils.Limit(w.Amp.Val, ampLimits.min, ampLimits.max)

	w.Pan.ModAmp = utils.Limit(w.Pan.ModAmp, panLimits.min, panLimits.max)
	w.Pan.Val = utils.Limit(w.Pan.Val, panLimits.min, panLimits.max)

	w.Freq.ModAmp = utils.Limit(w.Freq.ModAmp, freqLimits.min, freqLimits.max)
	w.Freq.Val = utils.Limit(w.Freq.Val, freqLimits.min, freqLimits.max)
}

func (w *Wavetable) signalValue(t, amp, freq float64) float64 {
	idx := int(math.Floor(t * float64(len(w.Table)) * freq))
	var val float64

	if l := len(w.Table); l > 0 {
		val = w.Table[idx%l]
	}

	y := amp * val
	w.integral += y / config.Config.SampleRate

	return y
}
