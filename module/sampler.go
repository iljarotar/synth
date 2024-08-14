package module

import (
	"github.com/iljarotar/synth/utils"
)

type Sampler struct {
	Module
	Name            string   `yaml:"name"`
	Amp             Input    `yaml:"amp"`
	Pan             Input    `yaml:"pan"`
	Freq            Input    `yaml:"freq"`
	Filters         []string `yaml:"filters"`
	Inputs          []string `yaml:"inputs"`
	Envelope        string   `yaml:"envelope"`
	inputs          []filterInputs
	lastTriggeredAt float64
	limits
	sampleRate float64
}

func (s *Sampler) Initialize(sampleRate float64) {
	s.sampleRate = sampleRate
	s.limits = limits{min: 0, max: sampleRate}
	s.limitParams()
	s.inputs = make([]filterInputs, len(s.Filters))
	s.current = stereo(0, s.Pan.Val)
}

func (s *Sampler) Next(t float64, modMap ModulesMap, filtersMap FiltersMap, envelopesMap EnvelopesMap) {
	amp := modulate(s.Amp, ampLimits, modMap)
	pan := modulate(s.Pan, panLimits, modMap)
	freq := modulate(s.Freq, s.limits, modMap)

	cfg := filterConfig{
		filterNames: s.Filters,
		inputs:      s.inputs,
		FiltersMap:  filtersMap,
	}

	x := s.sample(t, freq, amp, modMap)
	y, newInputs := cfg.applyFilters(x)
	y = applyEnvelope(y, s.Envelope, envelopesMap)
	s.integral += y / s.sampleRate
	s.inputs = newInputs
	s.current = stereo(y, pan)
}

func (s *Sampler) sample(t, freq, amp float64, modMap ModulesMap) float64 {
	if freq == 0 {
		return s.current.Mono
	}
	secondsBetweenTwoBeats := 1 / freq
	if t-s.lastTriggeredAt >= secondsBetweenTwoBeats {
		s.lastTriggeredAt = t
		return s.getCurrentOutputValue(modMap) * amp
	}
	return s.current.Mono
}

func (s *Sampler) getCurrentOutputValue(modMap ModulesMap) float64 {
	var y float64
	for _, m := range s.Inputs {
		mod, ok := modMap[m]
		if !ok {
			continue
		}
		y += mod.Current().Mono
	}
	return y
}

func (s *Sampler) limitParams() {
	s.Amp.Val = utils.Limit(s.Amp.Val, ampLimits.min, ampLimits.max)
	s.Amp.ModAmp = utils.Limit(s.Amp.ModAmp, ampLimits.min, ampLimits.max)

	s.Pan.Val = utils.Limit(s.Pan.Val, panLimits.min, panLimits.max)
	s.Pan.ModAmp = utils.Limit(s.Pan.ModAmp, panLimits.min, panLimits.max)

	s.Freq.Val = utils.Limit(s.Freq.Val, s.limits.min, s.limits.max)
	s.Freq.ModAmp = utils.Limit(s.Freq.ModAmp, s.limits.min, s.limits.max)
}
