package module

import "github.com/iljarotar/synth/calc"

type (
	Delay struct {
		Module
		Time float64 `yaml:"time"`
		Gain float64 `yaml:"gain"`
		In   string  `yaml:"in"`
		CV   string  `yaml:"cv"`
		Mod  string  `yaml:"mod"`
		Fade float64 `yaml:"fade"`

		sampleRate float64
		c          *comb

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
	d.Gain = calc.Limit(d.Gain, combMixRange)
	d.Fade = calc.Limit(d.Fade, fadeRange)
	d.c = &comb{}
	d.c.initialize(d.Time/1000, sampleRate)

	d.mixFader = &fader{
		current: d.Gain,
		target:  d.Gain,
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
		d.mixFader.target = new.Gain
		d.mixFader.initialize(d.Fade, d.sampleRate)
	}
	if d.c != nil {
		d.c.update(d.Time / 1000)
	}
}

func (d *Delay) Step(modules *ModuleMap) {
	mix := d.Gain
	if d.CV != "" {
		mix = cv(combMixRange, getMono(modules, d.CV))
	}
	mix = modulate(mix, combMixRange, getMono(modules, d.Mod))

	val := d.c.step(getMono(modules, d.In), mix)
	d.current = Output{
		Mono:  val,
		Left:  val / 2,
		Right: val / 2,
	}

	d.fade()
}

func (d *Delay) fade() {
	if d.mixFader != nil {
		d.Gain = d.mixFader.fade()
	}
}
