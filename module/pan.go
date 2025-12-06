package module

import "github.com/iljarotar/synth/calc"

type (
	Pan struct {
		Module
		Pan float64 `yaml:"pan"`
		Mod string  `yaml:"mod"`
		In  string  `yaml:"in"`
	}

	PanMap map[string]*Pan
)

func (m PanMap) Initialize() {
	for _, p := range m {
		if p == nil {
			continue
		}
		p.initialize()
	}
}

func (p *Pan) initialize() {
	p.Pan = calc.Limit(p.Pan, panRange)
}

func (p *Pan) Update(new *Pan) {
	if new == nil {
		return
	}

	p.Pan = new.Pan
	p.Mod = new.Mod
	p.In = new.In
}

func (p *Pan) Step(modules ModuleMap) {
	pan := modulate(p.Pan, panRange, getMono(modules[p.Mod]))
	percent := calc.Percentage(pan, panRange)
	in := getMono(modules[p.In])

	p.current = Output{
		Mono:  in,
		Right: in * percent,
		Left:  in * (1 - percent),
	}
}
