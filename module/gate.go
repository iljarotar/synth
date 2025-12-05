package module

import (
	"math"

	"github.com/iljarotar/synth/calc"
)

type (
	Gate struct {
		Module
		BPM        float64   `yaml:"bpm"`
		CV         string    `yaml:"cv"`
		Mod        string    `yaml:"mod"`
		Signal     []float64 `yaml:"signal"`
		sampleRate float64
		idx        float64
	}

	GateMap map[string]*Gate
)

func (m GateMap) Initialze(gates GateMap, sampleRate float64) {
	for name, g := range m {
		if g == nil {
			continue
		}

		var gate *Gate
		if gt, ok := gates[name]; ok {
			gate = gt
		}

		g.initialze(gate, sampleRate)
	}
}

func (g *Gate) initialze(gate *Gate, sampleRate float64) {
	g.sampleRate = sampleRate
	g.BPM = calc.Limit(g.BPM, bpmRange)

	for i, val := range g.Signal {
		if val <= 0 {
			g.Signal[i] = -1
		} else {
			g.Signal[i] = 1
		}
	}

	if gate != nil {
		g.idx = gate.idx
	}
}

func (g *Gate) Step(modules ModuleMap) {
	length := len(g.Signal)
	val := g.Signal[int(math.Floor(g.idx))%length]
	g.current = Output{
		Mono:  val,
		Left:  val / 2,
		Right: val / 2,
	}

	bpm := g.BPM
	if g.CV != "" {
		bpm = cv(bpmRange, getMono(modules[g.CV]))
	}

	bpm = modulate(bpm, bpmRange, getMono(modules[g.Mod]))
	spb := samplesPerBeat(g.sampleRate, bpm)
	if spb == 0 {
		return
	}

	g.idx += 1 / spb
}

func samplesPerBeat(sampleRate, bpm float64) float64 {
	if bpm == 0 {
		return math.Inf(1)
	}

	return sampleRate * 60 / bpm
}
