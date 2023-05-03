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
	Noise              []*module.Noise      `yaml:"noise"`
	Custom             []*module.Custom     `yaml:"custom"`
	Phase              float64
	oscMap             module.OscillatorsMap
	noiseMap           module.NoiseMap
	customMap          module.CustomMap
	step, volumeMemory float64
	next               chan bool
}

func (s *Synth) Initialize() {
	s.step = 1 / config.Config.SampleRate
	s.Volume = utils.Limit(s.Volume, 0, 1)
	s.volumeMemory = s.Volume
	s.Volume = 0 // start muted
	s.next = make(chan bool)

	for _, osc := range s.Oscillators {
		osc.Initialize()
	}

	for _, n := range s.Noise {
		n.Initialize()
	}

	for _, c := range s.Custom {
		c.Initialize()
	}

	s.makeOscillatorsMap()
	s.makeNoiseMap()
	s.makeCustomMap()
}

func (s *Synth) Play(input chan<- struct{ Left, Right float32 }) {
	defer close(input)

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
		case next := <-s.next:
			if next {
				s.next <- true
			} else {
				return
			}
		default:
		}
	}
}

func (s *Synth) Stop() {
	s.next <- false
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

		noise, ok := s.noiseMap[o]
		if ok {
			left += noise.Current.Left
			right += noise.Current.Right
		}

		custom, ok := s.customMap[o]
		if ok {
			left += custom.Current.Left
			right += custom.Current.Right
		}
	}

	return left, right
}

func (s *Synth) updateCurrentValues() {
	for _, o := range s.Oscillators {
		osc := o
		osc.Next(s.Phase, s.oscMap, s.customMap)
	}

	for _, n := range s.Noise {
		n.Next(s.oscMap, s.customMap)
	}

	for _, c := range s.Custom {
		c.Next(s.Phase, s.oscMap, s.customMap)
	}

	s.Phase += s.step
}

func (s *Synth) makeOscillatorsMap() {
	oscMap := make(module.OscillatorsMap)

	for _, osc := range s.Oscillators {
		oscMap[osc.Name] = osc
	}

	s.oscMap = oscMap
}

func (s *Synth) makeNoiseMap() {
	noiseMap := make(module.NoiseMap)

	for _, noise := range s.Noise {
		noiseMap[noise.Name] = noise
	}

	s.noiseMap = noiseMap
}

func (s *Synth) makeCustomMap() {
	cMap := make(module.CustomMap)

	for _, custom := range s.Custom {
		cMap[custom.Name] = custom
	}

	s.customMap = cMap
}

func secondsToStep(seconds, delta float64) float64 {
	steps := math.Round(seconds * config.Config.SampleRate)
	step := 1 / steps
	return step
}
