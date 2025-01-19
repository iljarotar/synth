package module

import (
	"math"

	"github.com/iljarotar/synth/utils"
)

type Wavetable struct {
	Module
	Name       string    `yaml:"name"`
	Table      []float64 `yaml:"table"`
	Freq       Input     `yaml:"freq"`
	Amp        Input     `yaml:"amp"`
	Pan        Input     `yaml:"pan"`
	Filters    []string  `yaml:"filters"`
	Envelope   *Envelope `yaml:"envelope"`
	inputs     []filterInputs
	sampleRate float64
}

func (w *Wavetable) Initialize(sampleRate float64) {
	if w.Envelope != nil {
		w.Envelope.Initialize()
	}
	w.sampleRate = sampleRate
	w.limitParams()
	w.Table = utils.Normalize(w.Table, -1, 1)
	w.inputs = make([]filterInputs, len(w.Filters))

	y := w.signalValue(0, w.Amp.Val, w.Freq.Val)
	w.current = stereo(y, w.Pan.Val)
}

func (w *Wavetable) Next(t float64, modMap ModulesMap, filtersMap FiltersMap) {
	if w.Envelope != nil {
		w.Envelope.Next(t, modMap)
	}

	pan := modulate(w.Pan, panLimits, modMap)
	amp := modulate(w.Amp, ampLimits, modMap)
	offset := w.getOffset(modMap)

	cfg := filterConfig{
		filterNames: w.Filters,
		inputs:      w.inputs,
		FiltersMap:  filtersMap,
	}

	x := w.signalValue(t, amp, offset)
	y, newInputs := cfg.applyFilters(x)
	y = applyEnvelope(y, w.Envelope)
	w.integral += y / w.sampleRate
	w.inputs = newInputs
	w.current = stereo(y, pan)
}

func (w *Wavetable) getOffset(modMap ModulesMap) float64 {
	var y float64

	for _, m := range w.Freq.Mod {
		mod, ok := modMap[m]
		if ok {
			y += mod.Integral()
		}
	}

	return y * w.Freq.ModAmp
}

func (w *Wavetable) signalValue(t, amp, offset float64) float64 {
	length := len(w.Table)
	if length == 0 {
		return 0
	}

	idx := int(math.Floor((t*w.Freq.Val + offset) * float64(length)))
	var val float64

	val = w.Table[idx%length]

	y := amp * val

	return y
}

func (w *Wavetable) limitParams() {
	w.Amp.ModAmp = utils.Limit(w.Amp.ModAmp, ampLimits.min, ampLimits.max)
	w.Amp.Val = utils.Limit(w.Amp.Val, ampLimits.min, ampLimits.max)

	w.Pan.ModAmp = utils.Limit(w.Pan.ModAmp, panLimits.min, panLimits.max)
	w.Pan.Val = utils.Limit(w.Pan.Val, panLimits.min, panLimits.max)

	w.Freq.ModAmp = utils.Limit(w.Freq.ModAmp, freqLimits.min, freqLimits.max)
	w.Freq.Val = utils.Limit(w.Freq.Val, freqLimits.min, freqLimits.max)
}
