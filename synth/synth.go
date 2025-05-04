package synth

import (
	"github.com/iljarotar/synth/module"
	"github.com/iljarotar/synth/utils"
)

const (
	maxInitTime = 7200
	maxVolume   = 2
)

type Output struct {
	Left, Right, Mono, Time float64
}

type Synth struct {
	Volume      float64              `yaml:"vol"`
	Out         []string             `yaml:"out"`
	Oscillators []*module.Oscillator `yaml:"oscillators"`
	Noises      []*module.Noise      `yaml:"noises"`
	Wavetables  []*module.Wavetable  `yaml:"wavetables"`
	Samplers    []*module.Sampler    `yaml:"samplers"`
	Sequences   []*module.Sequence   `yaml:"sequences"`
	Filters     []*module.Filter     `yaml:"filters"`
	Time        float64              `yaml:"time"`

	VolumeMemory         float64
	sampleRate           float64
	modMap               module.ModulesMap
	filtersMap           module.FiltersMap
	timeStep, volumeStep float64
	notifyFadeoutChan    chan<- bool
}

func (s *Synth) Initialize(sampleRate float64) error {
	s.timeStep = 1 / sampleRate
	s.sampleRate = sampleRate
	s.Volume = utils.Limit(s.Volume, 0, maxVolume)
	s.Time = utils.Limit(s.Time, 0, maxInitTime)
	s.VolumeMemory = s.Volume
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

func (s *Synth) Next() Output {
	left, right, mono := s.getCurrentValue()
	s.adjustVolume()

	left *= s.Volume
	right *= s.Volume
	mono *= s.Volume

	return Output{
		Left:  left,
		Right: right,
		Mono:  mono,
		Time:  s.Time,
	}
}

func (s *Synth) SetVolume(volume float64) {
	vol := min(volume, maxVolume)
	vol = max(vol, 0)
	s.VolumeMemory = vol
	s.Volume = vol
}

func (s *Synth) FadeIn(duration float64) {
	s.volumeStep = secondsToStep(duration, s.VolumeMemory-s.Volume, s.sampleRate)
}

func (s *Synth) FadeOut(duration float64) {
	s.volumeStep = secondsToStep(duration, -s.Volume, s.sampleRate)
}

func (s *Synth) NotifyFadeout(done chan<- bool) {
	s.notifyFadeoutChan = done
}

func (s *Synth) adjustVolume() {
	if s.volumeStep == 0 {
		if s.notifyFadeoutChan != nil {
			s.notifyFadeoutChan <- true
			close(s.notifyFadeoutChan)
			s.notifyFadeoutChan = nil
		}
		return
	}
	s.Volume += s.volumeStep

	if s.volumeStep > 0 && s.Volume >= s.VolumeMemory {
		s.volumeStep = 0
		s.Volume = s.VolumeMemory
		return
	}

	if s.Volume <= 0 {
		s.volumeStep = 0
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
		osc.Next(s.Time, s.modMap, s.filtersMap)
	}

	for _, n := range s.Noises {
		n.Next(s.Time, s.modMap, s.filtersMap)
	}

	for _, c := range s.Wavetables {
		c.Next(s.Time, s.modMap, s.filtersMap)
	}

	for _, smplr := range s.Samplers {
		smplr.Next(s.Time, s.modMap, s.filtersMap)
	}

	for _, sq := range s.Sequences {
		sq.Next(s.Time, s.modMap, s.filtersMap)
	}

	for _, f := range s.Filters {
		f.NextCoeffs(s.modMap)
	}

	s.Time += s.timeStep
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
