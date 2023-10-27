package module

import "github.com/iljarotar/synth/utils"

type Envelope struct {
	Module
	Name         string   `yaml:"name"`
	Attack       Param    `yaml:"attack"`
	Decay        Param    `yaml:"decay"`
	Sustain      Param    `yaml:"sustain"`
	Release      Param    `yaml:"release"`
	Peak         Param    `yaml:"peak"`
	SustainLevel Param    `yaml:"sustain-level"`
	Trigger      []string `yaml:"trigger"`
	Threshold    Param    `yaml:"threshold"`
	lastInput    float64  // TODO: use to check if from last to current input the threshold was exceeded; if so -> trigger
	triggeredAt  float64  // TODO: keep track of the last time the envelope was triggered
}

type envelopeConfig struct {
	Attack       float64
	Decay        float64
	Sustain      float64
	Release      float64
	Peak         float64
	SustainLevel float64
	Threshold    float64
}

func (e *Envelope) Initialize() {
	e.limitParams()
	e.current = output{Mono: 0, Left: 0, Right: 0}
}

func (e *Envelope) Next(t float64, modMap ModulesMap) {
	attack := utils.Limit(e.Attack.Val+modulate(e.Attack.Mod, modMap)*e.Attack.ModAmp, envelopeLimits.min, envelopeLimits.max)
	decay := utils.Limit(e.Decay.Val+modulate(e.Decay.Mod, modMap)*e.Decay.ModAmp, envelopeLimits.min, envelopeLimits.max)
	sustain := utils.Limit(e.Sustain.Val+modulate(e.Sustain.Mod, modMap)*e.Sustain.ModAmp, envelopeLimits.min, envelopeLimits.max)
	release := utils.Limit(e.Release.Val+modulate(e.Release.Mod, modMap)*e.Release.ModAmp, envelopeLimits.min, envelopeLimits.max)

	peak := utils.Limit(e.Peak.Val+modulate(e.Peak.Mod, modMap)*e.Peak.ModAmp, ampLimits.min, ampLimits.max)
	sustainLevel := utils.Limit(e.SustainLevel.Val+modulate(e.SustainLevel.Mod, modMap)*e.SustainLevel.ModAmp, ampLimits.min, ampLimits.max)
	threshold := utils.Limit(e.Threshold.Val+modulate(e.Threshold.Mod, modMap)*e.Threshold.ModAmp, ampLimits.min, ampLimits.max)

	envelope := envelopeConfig{
		Attack:       attack,
		Decay:        decay,
		Sustain:      sustain,
		Release:      release,
		Peak:         peak,
		SustainLevel: sustainLevel,
		Threshold:    threshold,
	}

	y := getCurrentValue(t, envelope)
	e.current = output{Mono: y, Left: 0, Right: 0}
}

func (e *Envelope) limitParams() {
	e.Attack.Val = utils.Limit(e.Attack.Val, envelopeLimits.min, envelopeLimits.max)
	e.Attack.ModAmp = utils.Limit(e.Attack.ModAmp, modLimits.min, modLimits.max)

	e.Decay.Val = utils.Limit(e.Decay.Val, envelopeLimits.min, envelopeLimits.max)
	e.Decay.ModAmp = utils.Limit(e.Decay.ModAmp, modLimits.min, modLimits.max)

	e.Sustain.Val = utils.Limit(e.Sustain.Val, envelopeLimits.min, envelopeLimits.max)
	e.Sustain.ModAmp = utils.Limit(e.Sustain.ModAmp, modLimits.min, modLimits.max)

	e.Release.Val = utils.Limit(e.Release.Val, envelopeLimits.min, envelopeLimits.max)
	e.Release.ModAmp = utils.Limit(e.Release.ModAmp, modLimits.min, modLimits.max)

	e.Peak.Val = utils.Limit(e.Peak.Val, ampLimits.min, ampLimits.max)
	e.Peak.ModAmp = utils.Limit(e.Peak.ModAmp, modLimits.min, modLimits.max)

	e.SustainLevel.Val = utils.Limit(e.SustainLevel.Val, ampLimits.min, ampLimits.max)
	e.SustainLevel.ModAmp = utils.Limit(e.SustainLevel.ModAmp, modLimits.min, modLimits.max)

	e.Threshold.Val = utils.Limit(e.Threshold.Val, ampLimits.min, ampLimits.max)
	e.Threshold.ModAmp = utils.Limit(e.Threshold.ModAmp, modLimits.min, modLimits.max)
}

func getCurrentValue(t float64, envelope envelopeConfig) float64 {
	// TODO: implement
	return 0
}
