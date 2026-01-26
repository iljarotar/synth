package module

import "github.com/iljarotar/synth/calc"

type (
	Pan struct {
		Module
		Pan  float64 `yaml:"pan"`
		Mod  string  `yaml:"mod"`
		In   string  `yaml:"in"`
		Fade float64 `yaml:"fade"`

		sampleRate float64

		panFader *fader
	}

	PanMap map[string]*Pan
)

func (m PanMap) Initialize(sampleRate float64) {
	for _, p := range m {
		if p == nil {
			continue
		}
		p.initialize(sampleRate)
	}
}

func (p *Pan) initialize(sampleRate float64) {
	p.sampleRate = sampleRate
	p.Pan = calc.Limit(p.Pan, panRange)
	p.Fade = calc.Limit(p.Fade, fadeRange)

	p.panFader = &fader{
		current: p.Pan,
		target:  p.Pan,
	}
	p.panFader.initialize(p.Fade, sampleRate)
}

func (p *Pan) Update(new *Pan) {
	if new == nil {
		return
	}

	p.Mod = new.Mod
	p.In = new.In
	p.Fade = new.Fade

	if p.panFader != nil {
		p.panFader.target = new.Pan
		p.panFader.initialize(p.Fade, p.sampleRate)
	}
}

func (p *Pan) Step(modules *ModuleMap) {
	pan := modulate(p.Pan, panRange, getMono(modules, p.Mod))
	percent := calc.Percentage(pan, panRange)
	in := getMono(modules, p.In)

	p.current = Output{
		Mono:  in,
		Right: in * percent,
		Left:  in * (1 - percent),
	}

	p.fade()
}

func (p *Pan) fade() {
	if p.panFader != nil {
		p.Pan = p.panFader.fade()
	}
}
