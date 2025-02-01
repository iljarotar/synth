package synth

import (
	"github.com/iljarotar/synth/module"
	"github.com/iljarotar/synth/utils"
)

const (
	maxInitTime = 7200
	maxVolume   = 2
)

type FadeDirection string

const (
	FadeDirectionIn  FadeDirection = "in"
	FadeDirectionOut FadeDirection = "out"
)

type Output struct {
	Left, Right, Mono, Time float64
}

type Synth struct {
	Volume             float64              `yaml:"vol"`
	Out                []string             `yaml:"out"`
	Oscillators        []*module.Oscillator `yaml:"oscillators"`
	Noises             []*module.Noise      `yaml:"noises"`
	Wavetables         []*module.Wavetable  `yaml:"wavetables"`
	Samplers           []*module.Sampler    `yaml:"samplers"`
	Sequences          []*module.Sequence   `yaml:"sequences"`
	Filters            []*module.Filter     `yaml:"filters"`
	time               float64
	sampleRate         float64
	modMap             module.ModulesMap
	filtersMap         module.FiltersMap
	step, volumeMemory float64
	fadeOutDone        chan bool
	fadeDirection      FadeDirection
	fadeDuration       float64
	active             bool
}

func (s *Synth) IncreaseVolume() {
	vol := s.Volume + 0.02
	if vol > maxVolume {
		vol = maxVolume
	}
	s.volumeMemory = vol
	s.Volume = vol
}

func (s *Synth) DecreaseVolume() {
	vol := s.Volume - 0.02
	if vol < 0 {
		vol = 0
	}
	s.volumeMemory = vol
	s.Volume = vol
}

func (s *Synth) initialize(sampleRate float64) error {
	s.step = 1 / sampleRate
	s.sampleRate = sampleRate
	s.Volume = utils.Limit(s.Volume, 0, maxVolume)
	s.time = utils.Limit(s.time, 0, maxInitTime)
	s.volumeMemory = s.Volume
	s.Volume = 0 // start muted

	for _, osc := range s.Oscillators {
		err := osc.Initialize(sampleRate)
		if err != nil {
			return err
		}
	}

	for _, n := range s.Noises {
		n.Initialize(sampleRate)
	}

	for _, c := range s.Wavetables {
		c.Initialize(sampleRate)
	}

	for _, smplr := range s.Samplers {
		smplr.Initialize(sampleRate)
	}

	for _, sq := range s.Sequences {
		err := sq.Initialize(sampleRate)
		if err != nil {
			return err
		}
	}

	for _, f := range s.Filters {
		f.Initialize(sampleRate)
	}
	s.makeMaps()

	return nil
}

func (s *Synth) play(outputChan chan<- Output) {
	defer close(outputChan)

	for s.active {
		left, right, mono := s.getCurrentValue()
		s.fade()
		left *= s.Volume
		right *= s.Volume
		mono *= s.Volume

		outputChan <- Output{
			Left:  left,
			Right: right,
			Mono:  mono,
			Time:  s.time,
		}
	}
}

func (s *Synth) startFading(direction FadeDirection, seconds float64) {
	s.fadeDirection = direction
	s.fadeDuration = seconds
}

func (s *Synth) notifyFadeOutDone(notify chan bool) {
	s.fadeOutDone = notify
}

func (s *Synth) fade() {
	if s.fadeDirection == FadeDirectionIn {
		s.fadeIn()
	} else {
		s.fadeOut()
	}
}

func (s *Synth) fadeIn() {
	if s.Volume == s.volumeMemory {
		return
	}

	step := secondsToStep(s.fadeDuration, s.volumeMemory-s.Volume, s.sampleRate)
	s.Volume += step
	s.fadeDuration -= 1 / s.sampleRate

	if s.Volume > s.volumeMemory {
		s.Volume = s.volumeMemory
	}
}

func (s *Synth) fadeOut() {
	if s.Volume == 0 {
		if s.fadeOutDone != nil {
			s.fadeOutDone <- true
			close(s.fadeOutDone)
			s.fadeOutDone = nil
		}
		return
	}

	step := secondsToStep(s.fadeDuration, s.Volume, s.sampleRate)
	s.Volume -= step
	s.fadeDuration -= 1 / s.sampleRate

	if s.Volume < 0 {
		s.Volume = 0
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
		osc.Next(s.time, s.modMap, s.filtersMap)
	}

	for _, n := range s.Noises {
		n.Next(s.time, s.modMap, s.filtersMap)
	}

	for _, c := range s.Wavetables {
		c.Next(s.time, s.modMap, s.filtersMap)
	}

	for _, smplr := range s.Samplers {
		smplr.Next(s.time, s.modMap, s.filtersMap)
	}

	for _, sq := range s.Sequences {
		sq.Next(s.time, s.modMap, s.filtersMap)
	}

	for _, f := range s.Filters {
		f.NextCoeffs(s.modMap)
	}

	s.time += s.step
}

func (s *Synth) makeMaps() {
	modMap := make(module.ModulesMap)
	filtersMap := make(module.FiltersMap)

	for _, osc := range s.Oscillators {
		modMap[osc.Name] = osc
	}

	for _, n := range s.Noises {
		modMap[n.Name] = n
	}

	for _, c := range s.Wavetables {
		modMap[c.Name] = c
	}

	for _, smplr := range s.Samplers {
		modMap[smplr.Name] = smplr
	}

	for _, sq := range s.Sequences {
		modMap[sq.Name] = sq
	}

	for _, f := range s.Filters {
		filtersMap[f.Name] = f
	}

	s.modMap = modMap
	s.filtersMap = filtersMap
}

func secondsToStep(seconds, delta, sampleRate float64) float64 {
	if seconds == 0 {
		return delta
	}
	steps := seconds * sampleRate
	step := delta / steps
	return step
}
