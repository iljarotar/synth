package synth

import (
	"github.com/iljarotar/synth/log"
	"github.com/iljarotar/synth/module"
	"github.com/iljarotar/synth/utils"
)

const (
	maxVolume = 1
)

type Output struct {
	Left, Right, Mono, Time float64
}

type Synth struct {
	Out    string  `yaml:"out"`
	Volume float64 `yaml:"vol"`

	Mixers            module.MixerMap      `yaml:"mixers"`
	Oscillators       module.OscillatorMap `yaml:"oscillators"`
	Time              float64
	VolumeMemory      float64
	sampleRate        float64
	volumeStep        float64
	notifyFadeoutChan chan<- bool
	logger            *log.Logger
	modules           module.ModulesMap
}

func (s *Synth) Initialize(sampleRate float64) error {
	s.sampleRate = sampleRate
	s.Volume = utils.Limit(s.Volume, 0, maxVolume)
	s.VolumeMemory = s.Volume
	s.Volume = 0
	s.makeModulesMap()

	err := s.Mixers.Initialize(sampleRate)
	if err != nil {
		return err
	}

	err = s.Oscillators.Initialize(sampleRate)
	if err != nil {
		return err
	}

	return nil
}

func (s *Synth) GetOutput() Output {
	s.step()
	s.adjustVolume()
	out := Output{Time: s.Time}

	if mod, ok := s.modules[s.Out]; ok {
		out.Left = mod.Current().Left * s.Volume
		out.Right = mod.Current().Right * s.Volume
		out.Mono = mod.Current().Mono * s.Volume
	}

	return out
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

func (s *Synth) step() {
	for _, m := range s.Mixers {
		m.Step(s.modules)
	}

	for _, osc := range s.Oscillators {
		osc.Step(s.Time)
	}

	s.Time += 1 / s.sampleRate
}

func secondsToStep(seconds, delta, sampleRate float64) float64 {
	if seconds == 0 {
		return delta
	}
	steps := seconds * sampleRate
	step := delta / steps
	return step
}

func (s *Synth) makeModulesMap() {
	s.modules = module.ModulesMap{}

	for name, m := range s.Mixers {
		s.modules[name] = m
	}
	for name, osc := range s.Oscillators {
		s.modules[name] = osc
	}
}
