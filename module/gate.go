package module

import (
	"math"

	"github.com/iljarotar/synth/calc"
)

type (
	Gate struct {
		Module
		Sub    int       `yaml:"sub"`
		Signal []float64 `yaml:"signal"`
		Index  int       `yaml:"index"`

		idx            float64
		samplesPerBeat float64
	}

	GateMap map[string]*Gate
)

func (m GateMap) Initialize(sampleRate, bpm float64) {
	for _, g := range m {
		if g == nil {
			continue
		}
		g.initialize(sampleRate, bpm)
	}
}

func (g *Gate) initialize(sampleRate, bpm float64) {
	g.samplesPerBeat = samplesPerBeat(sampleRate, bpm, g.Sub)
	g.Sub = int(calc.Limit(float64(g.Sub), subdivisionsRange))
	g.Index = int(calc.Limit(float64(g.Index), calc.Range{Min: 0, Max: float64(len(g.Signal) - 1)}))
	g.idx = float64(g.Index)

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
	g.Signal = new.Signal
	g.Sub = new.Sub
	g.samplesPerBeat = new.samplesPerBeat
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

	g.idx += 1 / g.samplesPerBeat
}

func samplesPerBeat(sampleRate, bpm float64, subdivisions int) float64 {
	if bpm == 0 || subdivisions == 0 {
		return math.Inf(1)
	}

	return sampleRate * 60 / (bpm * float64(subdivisions))
}
