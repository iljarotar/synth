package module

import (
	"math"

	"github.com/iljarotar/synth/utils"
)

type Envelope struct {
	Module
	Name            string  `yaml:"name"`
	Attack          Param   `yaml:"attack"`
	Decay           Param   `yaml:"decay"`
	Sustain         Param   `yaml:"sustain"`
	Release         Param   `yaml:"release"`
	Peak            Param   `yaml:"peak"`
	SustainLevel    Param   `yaml:"sustain-level"`
	TimeShift       float64 `yaml:"time-shift"`
	BPM             Param   `yaml:"bpm"`
	lastTriggeredAt *float64
	currentConfig   envelopeConfig
}

type envelopeConfig struct {
	attack       float64
	decay        float64
	sustain      float64
	release      float64
	peak         float64
	sustainLevel float64
}

func (e *Envelope) Initialize() {
	e.limitParams()
	e.current = output{Mono: 0, Left: 0, Right: 0}
}

func (e *Envelope) Next(t float64, modMap ModulesMap) {
	bpm := utils.Limit(e.BPM.Val+modulate(e.BPM.Mod, modMap)*e.BPM.ModAmp, bpmLimits.min, bpmLimits.max)
	e.trigger(t, bpm, modMap)
	y := e.getCurrentValue(t)
	e.current = output{Mono: y, Left: 0, Right: 0}
}

func (e *Envelope) getCurrentConfig(t float64, modMap ModulesMap) {
	attack := utils.Limit(e.Attack.Val+modulate(e.Attack.Mod, modMap)*e.Attack.ModAmp, envelopeLimits.min, envelopeLimits.max)
	decay := utils.Limit(e.Decay.Val+modulate(e.Decay.Mod, modMap)*e.Decay.ModAmp, envelopeLimits.min, envelopeLimits.max)
	sustain := utils.Limit(e.Sustain.Val+modulate(e.Sustain.Mod, modMap)*e.Sustain.ModAmp, envelopeLimits.min, envelopeLimits.max)
	release := utils.Limit(e.Release.Val+modulate(e.Release.Mod, modMap)*e.Release.ModAmp, envelopeLimits.min, envelopeLimits.max)

	peak := utils.Limit(e.Peak.Val+modulate(e.Peak.Mod, modMap)*e.Peak.ModAmp, ampLimits.min, ampLimits.max)
	sustainLevel := utils.Limit(e.SustainLevel.Val+modulate(e.SustainLevel.Mod, modMap)*e.SustainLevel.ModAmp, ampLimits.min, ampLimits.max)

	config := envelopeConfig{
		attack:       attack,
		decay:        decay,
		sustain:      sustain,
		release:      release,
		peak:         peak,
		sustainLevel: sustainLevel,
	}
	e.currentConfig = config
}

func (e *Envelope) trigger(t, bpm float64, modMap ModulesMap) {
	if bpm == 0 {
		return
	}
	secondsBetweenTwoBeats := 60 / bpm
	var triggerAt float64
	if t >= e.TimeShift {
		numberOfTriggersMinusOne := math.Floor((t - e.TimeShift) / secondsBetweenTwoBeats)
		triggerAt = numberOfTriggersMinusOne*secondsBetweenTwoBeats + e.TimeShift
	} else {
		numberOfTriggers := math.Ceil((e.TimeShift - t) / secondsBetweenTwoBeats)
		triggerAt = e.TimeShift - numberOfTriggers*secondsBetweenTwoBeats
	}

	oldLastTriggeredAt := e.lastTriggeredAt
	if oldLastTriggeredAt == nil {
		e.lastTriggeredAt = &triggerAt
		e.getCurrentConfig(t, modMap)
		return
	}

	if t-*e.lastTriggeredAt >= secondsBetweenTwoBeats {
		newLastTriggeredAt := t
		e.lastTriggeredAt = &newLastTriggeredAt
		e.getCurrentConfig(t, modMap)
	}
}

func (e *Envelope) getCurrentValue(t float64) float64 {
	if e.lastTriggeredAt == nil {
		return 0
	}
	attackEnd := e.currentConfig.attack
	decayEnd := attackEnd + e.currentConfig.decay
	sustainEnd := decayEnd + e.currentConfig.sustain
	releaseEnd := sustainEnd + e.currentConfig.release
	timeSinceLastTrigger := t - *e.lastTriggeredAt

	switch {
	case timeSinceLastTrigger <= attackEnd:
		return attackFunc(e.currentConfig, *e.lastTriggeredAt)(t)
	case timeSinceLastTrigger <= decayEnd:
		return decayFunc(e.currentConfig, *e.lastTriggeredAt)(t)
	case timeSinceLastTrigger <= sustainEnd:
		return e.currentConfig.sustainLevel
	case timeSinceLastTrigger <= releaseEnd:
		return releaseFunc(e.currentConfig, *e.lastTriggeredAt)(t)
	default:
		return 0
	}
}

type stageFunc func(t float64) float64

func attackFunc(envelope envelopeConfig, triggeredAt float64) stageFunc {
	m := envelope.peak / envelope.attack
	c := -m * triggeredAt

	return func(t float64) float64 {
		return m*t + c
	}
}

func decayFunc(envelope envelopeConfig, triggeredAt float64) stageFunc {
	m := -(envelope.peak - envelope.sustainLevel) / envelope.decay
	c := -m*(triggeredAt+envelope.attack) + envelope.peak

	return func(t float64) float64 {
		return m*t + c
	}
}

func releaseFunc(envelope envelopeConfig, triggeredAt float64) stageFunc {
	m := -envelope.sustainLevel / envelope.release
	c := -m*(triggeredAt+envelope.attack+envelope.decay+envelope.sustain) + envelope.sustainLevel

	return func(t float64) float64 {
		return m*t + c
	}
}

func (e *Envelope) limitParams() {
	e.Attack.Val = utils.Limit(e.Attack.Val, envelopeLimits.min, envelopeLimits.max)
	e.Attack.ModAmp = utils.Limit(e.Attack.ModAmp, envelopeLimits.min, envelopeLimits.max)

	e.Decay.Val = utils.Limit(e.Decay.Val, envelopeLimits.min, envelopeLimits.max)
	e.Decay.ModAmp = utils.Limit(e.Decay.ModAmp, envelopeLimits.min, envelopeLimits.max)

	e.Sustain.Val = utils.Limit(e.Sustain.Val, envelopeLimits.min, envelopeLimits.max)
	e.Sustain.ModAmp = utils.Limit(e.Sustain.ModAmp, envelopeLimits.min, envelopeLimits.max)

	e.Release.Val = utils.Limit(e.Release.Val, envelopeLimits.min, envelopeLimits.max)
	e.Release.ModAmp = utils.Limit(e.Release.ModAmp, envelopeLimits.min, envelopeLimits.max)

	e.Peak.Val = utils.Limit(e.Peak.Val, ampLimits.min, ampLimits.max)
	e.Peak.ModAmp = utils.Limit(e.Peak.ModAmp, ampLimits.min, ampLimits.max)

	e.SustainLevel.Val = utils.Limit(e.SustainLevel.Val, ampLimits.min, ampLimits.max)
	e.SustainLevel.ModAmp = utils.Limit(e.SustainLevel.ModAmp, ampLimits.min, ampLimits.max)

	e.BPM.Val = utils.Limit(e.BPM.Val, bpmLimits.min, bpmLimits.max)
	e.BPM.ModAmp = utils.Limit(e.BPM.ModAmp, bpmLimits.min, bpmLimits.max)
}
