package module

import (
	"github.com/iljarotar/synth/calc"
)

type (
	Envelope struct {
		Module
		Attack      float64 `yaml:"attack"`
		Decay       float64 `yaml:"decay"`
		Release     float64 `yaml:"release"`
		Peak        float64 `yaml:"peak"`
		Level       float64 `yaml:"level"`
		Gate        string  `yaml:"gate"`
		triggeredAt float64
		releasedAt  float64
		gateValue   float64
		level       float64
	}

	EnvelopeMap map[string]*Envelope
)

func (m EnvelopeMap) Initialize() {
	for _, e := range m {
		if e == nil {
			continue
		}
		e.initialize()
	}
}

func (e *Envelope) initialize() {
	e.Attack = calc.Limit(e.Attack, envelopeRange)
	e.Decay = calc.Limit(e.Decay, envelopeRange)
	e.Release = calc.Limit(e.Release, envelopeRange)
	e.Peak = calc.Limit(e.Peak, gainRange)
	e.Level = calc.Limit(e.Level, gainRange)
}

func (e *Envelope) Update(new *Envelope) {
	if new == nil {
		return
	}

	e.Attack = new.Attack
	e.Decay = new.Decay
	e.Release = new.Release
	e.Peak = new.Peak
	e.Level = new.Level
	e.Gate = new.Gate
}

func (e *Envelope) Step(t float64, modules ModuleMap) {
	gateValue := getMono(modules[e.Gate])

	switch {
	case e.gateValue <= 0 && gateValue > 0:
		e.triggeredAt = t
	case e.gateValue > 0 && gateValue <= 0:
		e.releasedAt = t
		e.level = calc.Transpose(e.current.Mono, outputRange, gainRange)
	default:
		// noop
	}

	val := calc.Transpose(e.getValue(t), gainRange, outputRange)
	e.current = Output{
		Mono:  val,
		Left:  val / 2,
		Right: val / 2,
	}

	e.gateValue = gateValue
}

func (e *Envelope) getValue(t float64) float64 {
	if e.releasedAt >= e.triggeredAt {
		if t-e.releasedAt > e.Release {
			return 0
		}
		return e.release(t)
	}

	switch {
	case t-e.triggeredAt < e.Attack:
		return e.attack(t)
	case t-e.triggeredAt < e.Attack+e.Decay:
		return e.decay(t)
	default:
		return e.Level
	}
}

func (e *Envelope) attack(t float64) float64 {
	start := e.triggeredAt
	end := start + e.Attack
	return linear(start, end, 0, e.Peak, t)
}

func (e *Envelope) decay(t float64) float64 {
	start := e.triggeredAt + e.Attack
	end := start + e.Decay
	return linear(start, end, e.Peak, e.Level, t)
}

func (e *Envelope) release(t float64) float64 {
	start := e.releasedAt
	end := start + e.Release
	return linear(start, end, e.level, 0, t)
}

func linear(startAt, endAt, startValue, targetValue, t float64) float64 {
	delta := endAt - startAt
	if delta == 0 {
		return targetValue
	}
	return ((targetValue-startValue)*(t-startAt) + startValue*delta) / delta
}
