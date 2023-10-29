package module

import (
	"math"

	"github.com/iljarotar/synth/utils"
)

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
	Negative     bool     `yaml:"negative"`
	lastInput    float64
	triggeredAt  float64
	triggered    bool
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

	e.checkTrigger(t, threshold, modMap)
	y := e.getCurrentValue(t, envelope)
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

func (e *Envelope) checkTrigger(t, threshold float64, modMap ModulesMap) {
	var sum float64
	for _, trigger := range e.Trigger {
		mod, ok := modMap[trigger]
		if ok {
			sum += mod.Current().Mono
		}
	}

	sum = math.Abs(sum)

	if e.lastInput < threshold && sum >= threshold {
		e.triggered = true
		e.triggeredAt = t
	}

	e.lastInput = sum
}

func (e *Envelope) getCurrentValue(t float64, envelope envelopeConfig) float64 {
	if !e.triggered {
		return 0
	}

	var f stageFunc
	attackEnd := e.triggeredAt + envelope.Attack
	decayEnd := attackEnd + envelope.Decay
	sustainEnd := decayEnd + envelope.Sustain
	releaseEnd := sustainEnd + envelope.Release

	switch {
	case t >= e.triggeredAt && t < attackEnd:
		f = attackFunc(envelope, e.triggeredAt)

	case t >= attackEnd && t < decayEnd:
		f = decayFunc(envelope, e.triggeredAt)

	case t >= decayEnd && t < sustainEnd:
		return envelope.SustainLevel

	case t >= sustainEnd && t < releaseEnd:
		f = releaseFunc(envelope, e.triggeredAt)

	default:
		e.triggered = false
		return 0
	}

	if e.Negative {
		return -f(t)
	}

	return f(t)
}

type stageFunc func(t float64) float64

func attackFunc(envelope envelopeConfig, triggeredAt float64) stageFunc {
	m := envelope.Peak / envelope.Attack
	c := -m * triggeredAt

	return func(t float64) float64 {
		return m*t + c
	}
}

func decayFunc(envelope envelopeConfig, triggeredAt float64) stageFunc {
	m := -(envelope.Peak - envelope.SustainLevel) / envelope.Decay
	c := -m*(triggeredAt+envelope.Attack) + envelope.Peak

	return func(t float64) float64 {
		return m*t + c
	}
}

func releaseFunc(envelope envelopeConfig, triggeredAt float64) stageFunc {
	m := -envelope.SustainLevel / envelope.Release
	c := -m*(triggeredAt+envelope.Attack+envelope.Decay+envelope.Sustain) + envelope.SustainLevel

	return func(t float64) float64 {
		return m*t + c
	}
}
