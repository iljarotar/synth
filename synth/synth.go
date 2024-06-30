package synth

import (
	"math"

	"github.com/iljarotar/synth/config"
	"github.com/iljarotar/synth/module"
	"github.com/iljarotar/synth/ui"
	"github.com/iljarotar/synth/utils"
)

type Synth struct {
	Volume             float64                 `yaml:"vol"`
	Out                []string                `yaml:"out"`
	Oscillators        []*module.Oscillator    `yaml:"oscillators"`
	Noises             []*module.Noise         `yaml:"noises"`
	CustomSignals      []*module.CustomSignal  `yaml:"custom-signals"`
	TextProcessors     []*module.TextProcessor `yaml:"text-processors"`
	Envelopes          []*module.Envelope      `yaml:"envelopes"`
	Time               float64                 `yaml:"time"`
	modMap             module.ModulesMap
	step, volumeMemory float64
	next               chan bool
}

func (s *Synth) Initialize() {
	s.step = 1 / config.Config.SampleRate
	s.Volume = utils.Limit(s.Volume, 0, 1)
	s.Time = utils.Limit(s.Time, 0, 7200)
	s.volumeMemory = s.Volume
	s.Volume = 0 // start muted
	s.next = make(chan bool)

	for _, osc := range s.Oscillators {
		osc.Initialize()
	}

	for _, n := range s.Noises {
		n.Initialize()
	}

	for _, c := range s.CustomSignals {
		c.Initialize()
	}

	for _, p := range s.TextProcessors {
		p.Initialize()
	}

	for _, e := range s.Envelopes {
		e.Initialize()
	}

	s.makeModulesMap()
}

func (s *Synth) Play(output chan<- struct{ Left, Right float32 }) {
	defer close(output)

	for {
		left, right, mono := s.getCurrentValue()
		left *= s.Volume
		right *= s.Volume
		mono *= s.Volume

		if mono > 1 && !ui.Logger.ShowingOverdriveWarning {
			ui.Logger.ShowOverdriveWarning(true)
		}

		y := struct{ Left, Right float32 }{Left: float32(left), Right: float32(right)}
		output <- y

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

func (s *Synth) getCurrentValue() (left, right, mono float64) {
	s.updateCurrentValues()
	left, right, mono = 0, 0, 0

	for _, o := range s.Out {
		mod, ok := s.modMap[o]
		if ok {
			left += mod.Current().Left
			right += mod.Current().Right
			mono += mod.Current().Mono
		}
	}

	return left, right, mono
}

func (s *Synth) updateCurrentValues() {
	for _, o := range s.Oscillators {
		osc := o
		osc.Next(s.Time, s.modMap)
	}

	for _, n := range s.Noises {
		n.Next(s.Time, s.modMap)
	}

	for _, c := range s.CustomSignals {
		c.Next(s.Time, s.modMap)
	}

	for _, p := range s.TextProcessors {
		p.Next(s.Time, s.modMap)
	}

	for _, e := range s.Envelopes {
		e.Next(s.Time, s.modMap)
	}

	s.Time += s.step
}

func (s *Synth) makeModulesMap() {
	modMap := make(module.ModulesMap)

	for _, osc := range s.Oscillators {
		modMap[osc.Name] = osc
	}

	for _, n := range s.Noises {
		modMap[n.Name] = n
	}

	for _, c := range s.CustomSignals {
		modMap[c.Name] = c
	}

	for _, p := range s.TextProcessors {
		modMap[p.Name] = p
	}

	for _, e := range s.Envelopes {
		modMap[e.Name] = e
	}

	s.modMap = modMap
}

func secondsToStep(seconds, delta float64) float64 {
	steps := math.Round(seconds * config.Config.SampleRate)
	step := 1 / steps
	return step
}
