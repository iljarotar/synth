package synth

import (
	"math"

	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/module"
	"github.com/iljarotar/synth/utils"
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
	next               chan bool
}

func (s *Synth) Initialize() {
	s.step = 1 / config.Instance.SampleRate
	s.Volume = utils.Limit(s.Volume, 0, 1)
	s.volumeMemory = s.Volume
	s.Volume = 0 // start muted
	s.next = make(chan bool)

	for _, osc := range s.Oscillators {
		osc.Initialize()
	}

	for _, f := range s.Filters {
		f.Initialize()
	}

	s.makeOscillatorsMap()
	s.makeFiltersMap()
}

func (s *Synth) Play(input chan<- struct{ Left, Right float32 }) {
	for {
		left, right := s.getCurrentValue()
		left *= s.Volume
		right *= s.Volume

		if l := len(s.Out); l > 0 {
			left /= float64(l)
			right /= float64(l)
		}

		y := struct{ Left, Right float32 }{Left: float32(left), Right: float32(right)}
		input <- y

		select {
		case <-s.next:
			s.next <- true
		default:
		}
	}
}

func (s *Synth) FadeOut(seconds float64) {
	step := secondsToStep(seconds, s.Volume)
	for s.Volume > 0 {
		s.Volume -= step
		s.next <- true
		<-s.next
	}

	if s.Volume < 0 {
		s.Volume = 0
	}
}

func (s *Synth) FadeIn(seconds float64) {
	step := secondsToStep(seconds, s.volumeMemory-s.Volume)
	for s.Volume < s.volumeMemory {
		s.Volume += step
		s.next <- true
		<-s.next
	}

	if s.Volume > s.volumeMemory {
		s.Volume = s.volumeMemory
	}
}

func (s *Synth) getCurrentValue() (left, right float64) {
	s.updateCurrentValues()
	left, right = 0, 0

	for _, o := range s.Out {
		osc, ok := s.oscMap[o]
		if ok {
			left += osc.Current.Left
			right += osc.Current.Right
		}
	}

	return left, right
}

func (s *Synth) updateCurrentValues() {
	for _, o := range s.Oscillators {
		osc := o
		osc.Next(s.oscMap, s.filtersMap, s.Phase)
	}

	for _, f := range s.Filters {
		f.Next(s.oscMap)
	}

	s.Phase += s.step
}

func (s *Synth) makeOscillatorsMap() {
	oscMap := make(module.Oscillators)

	for _, osc := range s.Oscillators {
		oscMap[osc.Name] = osc
	}

	s.oscMap = oscMap
}

func (s *Synth) makeFiltersMap() {
	filtersMap := make(module.Filters)

	for _, f := range s.Filters {
		filtersMap[f.Name] = f
	}

	s.filtersMap = filtersMap
}

func secondsToStep(seconds, delta float64) float64 {
	steps := math.Round(seconds * config.Instance.SampleRate)
	step := 1 / steps
	return step
}
