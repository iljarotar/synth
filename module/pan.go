package module

import "github.com/iljarotar/synth/calc"

type Pan struct {
	Module
	Pan float64 `yaml:"pan"`
	Mod string  `yaml:"mod"`
	In  string  `yaml:"in"`
}

type PanMap map[string]*Pan

func (m PanMap) Initialize() {
	for _, p := range m {
		p.initialize()
	}
}

func (p *Pan) initialize() {
	p.Pan = calc.Limit(p.Pan, panLimits)
}

func (p *Pan) Step(modules ModuleMap) {
	pan := modulate(p.Pan, panLimits, getMono(modules[p.Mod]))
	percent := calc.Percentage(pan, panLimits)
	in := getMono(modules[p.In])

	p.current = Output{
		Mono:  in,
		Right: in * percent,
		Left:  in * (1 - percent),
	}
}
