package synth

import (
	"time"

	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/module"
)

type Synth struct {
	Volume             float64              `yaml:"volume"`
	Out                []string             `yaml:"out"`
	Oscillators        []*module.Oscillator `yaml:"oscillators"`
	Filters            []*module.Filter     `yaml:"filters"`
	Phase              float64
	oscMap             module.Oscillators
	filtersMap         module.Filters
	step, volumeMemory float64
}

func (s *Synth) Initialize() {
	s.step = 1 / config.Instance.SampleRate()

	if s.Volume == 0 {
		s.Volume = 1
	}

	s.volumeMemory = s.Volume
	s.Volume = 0 // start muted

	for i := range s.Oscillators {
		s.Oscillators[i].Initialize()
	}

	for i := range s.Filters {
		s.Filters[i].Initialize()
	}

	s.makeOscillatorsMap()
	s.makeFiltersMap()
}

func (s *Synth) Play(input chan<- float32) {
	for {
		y := s.getCurrentValue() * s.Volume
		s.Phase += s.step

		if len(s.Out) > 0 {
			y /= float64(len(s.Out))
		}

		input <- float32(y)
	}
}

func (s *Synth) FadeOut(step float64) {
	sampleRate := config.Instance.SampleRate()
	for s.Volume > 0 {
		s.Volume -= step
		time.Sleep(time.Second / time.Duration(sampleRate))
	}
}

func (s *Synth) FadeIn(step float64) {
	sampleRate := config.Instance.SampleRate()
	for s.Volume < s.volumeMemory {
		s.Volume += step
		time.Sleep(time.Second / time.Duration(sampleRate))
	}
}

func (s *Synth) getCurrentValue() float64 {
	s.updateValues()
	var y float64

	for i := range s.Out {
		osc := s.oscMap[s.Out[i]]
		y += osc.Current
	}

	return y
}

func (s *Synth) updateValues() {
	s.updateFilters()

	for i := range s.Oscillators {
		osc := s.Oscillators[i]
		amp := getAmp(*osc, s.oscMap)
		shift := getPhase(*osc, s.oscMap)

		for j := range osc.Filters {
			f, ok := s.filtersMap[osc.Filters[j]]
			if ok {
				amp *= f.Apply(osc.Freq)
			}
		}

		osc.Current = osc.Signal(osc.Freq*(s.Phase+shift)) * amp
	}
}

func (s *Synth) updateFilters() {
	for i := range s.Filters {
		f := s.Filters[i]
		x := make([]float64, 0)

		for j := range f.Cutoff.Mod {
			osc, ok := s.oscMap[f.Cutoff.Mod[j]]
			if ok {
				x = append(x, osc.Current)
			}
		}

		f.UpdateCutoff(x...)
	}
}

func getAmp(osc module.Oscillator, oscMap module.Oscillators) float64 {
	amp := 1.0

	if osc.Amp != nil {
		amp = osc.Amp.Val

		for j := range osc.Amp.Mod {
			mod, ok := oscMap[osc.Amp.Mod[j]]
			if ok {
				amp += mod.Current
			}
		}
	}

	return amp
}

func getPhase(osc module.Oscillator, oscMap module.Oscillators) float64 {
	phase := osc.Phase.Val

	for j := range osc.Phase.Mod {
		mod, ok := oscMap[osc.Phase.Mod[j]]
		if ok {
			phase += mod.Current
		}
	}

	return phase
}

func (s *Synth) makeOscillatorsMap() {
	oscMap := make(module.Oscillators)

	for i := range s.Oscillators {
		osc := s.Oscillators[i]
		oscMap[osc.Name] = osc
	}

	s.oscMap = oscMap
}

func (s *Synth) makeFiltersMap() {
	filtersMap := make(module.Filters)

	for i := range s.Filters {
		f := s.Filters[i]
		filtersMap[f.Name] = f
	}

	s.filtersMap = filtersMap
}
