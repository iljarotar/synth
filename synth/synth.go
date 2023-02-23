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

func (s *Synth) Play(input chan<- struct{ Left, Right float32 }) {
	for {
		left, right := s.getCurrentValue()
		left *= s.Volume
		right *= s.Volume

		if len(s.Out) > 0 {
			left /= float64(len(s.Out))
			right /= float64(len(s.Out))
		}

		y := struct{ Left, Right float32 }{Left: float32(left), Right: float32(right)}
		input <- y
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

func (s *Synth) getCurrentValue() (left, right float64) {
	s.updateCurrentValues()
	left, right = 0, 0

	for i := range s.Out {
		osc, ok := s.oscMap[s.Out[i]]
		if ok {
			left += osc.Current.Left
			right += osc.Current.Right
		}
	}

	return left, right
}

func (s *Synth) updateCurrentValues() {
	for i := range s.Oscillators {
		osc := s.Oscillators[i]
		osc.Next(s.oscMap, s.filtersMap, s.Phase)
	}

	for i := range s.Filters {
		f := s.Filters[i]
		f.Next(s.oscMap)
	}

	s.Phase += s.step
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
