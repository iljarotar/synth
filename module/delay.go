package module

import "github.com/iljarotar/synth/calc"

type (
	Delay struct {
		Module
		Time float64 `yaml:"time"`
		Mix  float64 `yaml:"mix"`
		In   string  `yaml:"in"`
		CV   string  `yaml:"cv"`
		Mod  string  `yaml:"mod"`
		Fade float64 `yaml:"fade"`

		sampleRate float64
		comb       *comb

		mixFader *fader
	}

	DelayMap map[string]*Delay
)

func (m DelayMap) Initialize(sampleRate float64) {
	for _, d := range m {
		if d == nil {
			continue
		}
		d.initialize(sampleRate)
	}
}

func (d *Delay) initialize(sampleRate float64) {
	d.sampleRate = sampleRate
	d.Time = calc.Limit(d.Time, combTimeRange)
	d.Mix = calc.Limit(d.Mix, combMixRange)
	d.Fade = calc.Limit(d.Fade, fadeRange)
	d.comb = &comb{}
	d.comb.initialize(d.Time/1000, sampleRate)

	d.mixFader = &fader{
		current: d.Mix,
		target:  d.Mix,
	}
	d.mixFader.initialize(d.Fade, sampleRate)
}

func (d *Delay) Update(new *Delay) {
	if new == nil {
		return
	}

	d.In = new.In
	d.CV = new.CV
	d.Mod = new.Mod
	d.Time = new.Time
	d.Fade = new.Fade

	if d.mixFader != nil {
		d.mixFader.target = new.Mix
		d.mixFader.initialize(d.Fade, d.sampleRate)
	}
	d.comb.update(d.Time / 1000)
}

func (d *Delay) Step(modules *ModuleMap) {
	mix := d.Mix
	if d.CV != "" {
		mix = cv(combMixRange, getMono(modules, d.CV))
	}
	mix = modulate(mix, combMixRange, getMono(modules, d.Mod))

	val := d.comb.step(getMono(modules, d.In), mix)
	d.current = Output{
		Mono:  val,
		Left:  val / 2,
		Right: val / 2,
	}

	d.fade()
}

func (d *Delay) fade() {
	if d.mixFader != nil {
		d.Mix = d.mixFader.fade()
	}
}
