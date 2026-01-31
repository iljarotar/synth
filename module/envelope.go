package module

import (
	"github.com/iljarotar/synth/calc"
)

type (
	Envelope struct {
		Module
		Attack  float64 `yaml:"attack"`
		Decay   float64 `yaml:"decay"`
		Release float64 `yaml:"release"`
		Peak    float64 `yaml:"peak"`
		Level   float64 `yaml:"level"`
		Gate    string  `yaml:"gate"`
		Fade    float64 `yaml:"fade"`

		triggeredAt float64
		releasedAt  float64
		gateValue   float64
		level       float64
		sampleRate  float64

		attackFader  *fader
		decayFader   *fader
		releaseFader *fader
		peakFader    *fader
		levelFader   *fader
	}

	EnvelopeMap map[string]*Envelope
)

func (m EnvelopeMap) Initialize(sampleRate float64) {
	for _, e := range m {
		if e == nil {
			continue
		}
		e.initialize(sampleRate)
	}
}

func (e *Envelope) initialize(sampleRate float64) {
	e.sampleRate = sampleRate
	e.Attack = calc.Limit(e.Attack, envelopeRange)
	e.Decay = calc.Limit(e.Decay, envelopeRange)
	e.Release = calc.Limit(e.Release, envelopeRange)
	e.Peak = calc.Limit(e.Peak, gainRange)
	e.Level = calc.Limit(e.Level, gainRange)
	e.Fade = calc.Limit(e.Fade, fadeRange)

	e.attackFader = &fader{
		current: e.Attack,
		target:  e.Attack,
	}
	e.decayFader = &fader{
		current: e.Decay,
		target:  e.Decay,
	}
	e.releaseFader = &fader{
		current: e.Release,
		target:  e.Release,
	}
	e.peakFader = &fader{
		current: e.Peak,
		target:  e.Peak,
	}
	e.levelFader = &fader{
		current: e.Level,
		target:  e.Level,
	}
	e.initializeFaders()
}

func (e *Envelope) Update(new *Envelope) {
	if new == nil {
		return
	}

	e.Gate = new.Gate
	e.Fade = new.Fade

	if e.attackFader != nil {
		e.attackFader.target = new.Attack
	}
	if e.decayFader != nil {
		e.decayFader.target = new.Decay
	}
	if e.releaseFader != nil {
		e.releaseFader.target = new.Release
	}
	if e.peakFader != nil {
		e.peakFader.target = new.Peak
	}
	if e.levelFader != nil {
		e.levelFader.target = new.Level
	}
	e.initializeFaders()
}

func (e *Envelope) Step(t float64, modules *ModuleMap) {
	gateValue := getMono(modules, e.Gate)

	switch {
	case e.gateValue <= 0 && gateValue > 0:
		e.triggeredAt = t
	case e.gateValue > 0 && gateValue <= 0:
		e.releasedAt = t
		e.level = calc.Transpose(e.current.Mono, cvRange, gainRange)
	default:
		// noop
	}

	val := calc.Transpose(e.getValue(t), gainRange, cvRange)
	e.current = Output{
		Mono:  val,
		Left:  val / 2,
		Right: val / 2,
	}

	e.gateValue = gateValue
	e.fade()
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

func (e *Envelope) initializeFaders() {
	if e.attackFader != nil {
		e.attackFader.initialize(e.Fade, e.sampleRate)
	}
	if e.decayFader != nil {
		e.decayFader.initialize(e.Fade, e.sampleRate)
	}
	if e.releaseFader != nil {
		e.releaseFader.initialize(e.Fade, e.sampleRate)
	}
	if e.peakFader != nil {
		e.peakFader.initialize(e.Fade, e.sampleRate)
	}
	if e.levelFader != nil {
		e.levelFader.initialize(e.Fade, e.sampleRate)
	}
}

func (e *Envelope) fade() {
	if e.attackFader != nil {
		e.Attack = e.attackFader.fade()
	}
	if e.decayFader != nil {
		e.Decay = e.decayFader.fade()
	}
	if e.releaseFader != nil {
		e.Release = e.releaseFader.fade()
	}
	if e.peakFader != nil {
		e.Peak = e.peakFader.fade()
	}
	if e.levelFader != nil {
		e.Level = e.levelFader.fade()
	}
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
