package module

import (
	"math"

	"github.com/iljarotar/synth/calc"
)

type (
	Gate struct {
		Module
		BPM    float64   `yaml:"bpm"`
		CV     string    `yaml:"cv"`
		Mod    string    `yaml:"mod"`
		Signal []float64 `yaml:"signal"`
		Index  int       `yaml:"index"`
		Fade   float64   `yaml:"fade"`

		sampleRate float64
		idx        float64

		bpmFader *fader
	}

	GateMap map[string]*Gate
)

func (m GateMap) Initialize(sampleRate float64) {
	for _, g := range m {
		if g == nil {
			continue
		}
		g.initialize(sampleRate)
	}
}

func (g *Gate) initialize(sampleRate float64) {
	g.sampleRate = sampleRate
	g.BPM = calc.Limit(g.BPM, bpmRange)
	g.Fade = calc.Limit(g.Fade, fadeRange)
	g.Index = int(calc.Limit(float64(g.Index), calc.Range{Min: 0, Max: float64(len(g.Signal) - 1)}))
	g.idx = float64(g.Index)

	g.bpmFader = &fader{
		current: g.BPM,
		target:  g.BPM,
	}
	g.bpmFader.initialize(g.Fade, sampleRate)

	for i, val := range g.Signal {
		if val <= 0 {
			g.Signal[i] = -1
		} else {
			g.Signal[i] = 1
		}
	}
}

func (g *Gate) Update(new *Gate) {
	if new == nil {
		return
	}

	g.CV = new.CV
	g.Mod = new.Mod
	g.Signal = new.Signal
	g.Fade = new.Fade

	if g.bpmFader != nil {
		g.bpmFader.target = new.BPM
		g.bpmFader.initialize(g.Fade, g.sampleRate)
	}
}

func (g *Gate) Step(modules *ModuleMap) {
	if len(g.Signal) < 1 {
		return
	}

	val := g.Signal[int(math.Floor(g.idx))%len(g.Signal)]
	g.current = Output{
		Mono:  val,
		Left:  val / 2,
		Right: val / 2,
	}

	bpm := g.BPM
	if g.CV != "" {
		bpm = cv(bpmRange, getMono(modules, g.CV))
	}

	bpm = modulate(bpm, bpmRange, getMono(modules, g.Mod))
	spb := samplesPerBeat(g.sampleRate, bpm)
	if spb == 0 {
		return
	}

	g.idx += 1 / spb
	g.fade()
}

func (g *Gate) fade() {
	if g.bpmFader != nil {
		g.BPM = g.bpmFader.fade()
	}
}

func samplesPerBeat(sampleRate, bpm float64) float64 {
	if bpm == 0 {
		return math.Inf(1)
	}

	return sampleRate * 60 / bpm
}
