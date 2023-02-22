package synth

import (
	"time"

	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/module"
)

type Synth struct {
	Volume             float64              `yaml:"vol"`
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
	s.updateCurrentValues()
	var y float64

	for i := range s.Out {
		osc, ok := s.oscMap[s.Out[i]]
		if ok {
			y += osc.Current
		}
	}

	return y
}

func (s *Synth) updateCurrentValues() {
	for i := range s.Oscillators {
		osc := s.Oscillators[i]
		osc.Next(s.oscMap, s.filtersMap, s.Phase)
	}

	for i := range s.Filters {
		f := s.Filters[i]
		f.Next(s.oscMap, s.Phase)
	}
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
