package synth

import (
	"github.com/iljarotar/synth/calc"
	"github.com/iljarotar/synth/log"
	"github.com/iljarotar/synth/module"
	"github.com/samber/lo"
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

	Gates       module.GateMap       `yaml:"gates"`
	Mixers      module.MixerMap      `yaml:"mixers"`
	Noises      module.NoiseMap      `yaml:"noises"`
	Oscillators module.OscillatorMap `yaml:"oscillators"`
	Pans        module.PanMap        `yaml:"pans"`
	Samplers    module.SamplerMap    `yaml:"samplers"`
	Wavetables  module.WavetableMap  `yaml:"wavetables"`

	Time              float64
	VolumeMemory      float64
	sampleRate        float64
	volumeStep        float64
	notifyFadeoutChan chan<- bool
	logger            *log.Logger
	modules           module.ModuleMap

	gates       []*module.Gate
	mixers      []*module.Mixer
	noises      []*module.Noise
	oscillators []*module.Oscillator
	pans        []*module.Pan
	samplers    []*module.Sampler
	wavetables  []*module.Wavetable
}

func (s *Synth) Initialize(sampleRate float64) error {
	s.sampleRate = sampleRate
	s.Volume = calc.Limit(s.Volume, calc.Range{
		Min: 0,
		Max: maxVolume,
	})
	s.VolumeMemory = s.Volume
	s.Volume = 0
	s.makeModulesMap()
	s.flattenModules()

	s.Gates.Initialze(sampleRate)

	err := s.Mixers.Initialize(sampleRate)
	if err != nil {
		return err
	}

	err = s.Oscillators.Initialize(sampleRate)
	if err != nil {
		return err
	}

	s.Pans.Initialize()
	s.Wavetables.Initialize(sampleRate)

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
	for _, g := range s.gates {
		g.Step(s.modules)
	}
	for _, m := range s.mixers {
		m.Step(s.modules)
	}
	for _, n := range s.noises {
		n.Step()
	}
	for _, osc := range s.oscillators {
		osc.Step(s.modules)
	}
	for _, p := range s.pans {
		p.Step(s.modules)
	}
	for _, smplr := range s.samplers {
		smplr.Step(s.modules)
	}
	for _, w := range s.wavetables {
		w.Step(s.modules)
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
	s.modules = module.ModuleMap{}

	for name, g := range s.Gates {
		s.modules[name] = g
	}
	for name, m := range s.Mixers {
		s.modules[name] = m
	}
	for name, n := range s.Noises {
		s.modules[name] = n
	}
	for name, osc := range s.Oscillators {
		s.modules[name] = osc
	}
	for name, p := range s.Pans {
		s.modules[name] = p
	}
	for name, smplr := range s.Samplers {
		s.modules[name] = smplr
	}
	for name, w := range s.Wavetables {
		s.modules[name] = w
	}
}

func (s *Synth) flattenModules() {
	s.gates = lo.Values(s.Gates)
	s.mixers = lo.Values(s.Mixers)
	s.noises = lo.Values(s.Noises)
	s.oscillators = lo.Values(s.Oscillators)
	s.pans = lo.Values(s.Pans)
	s.samplers = lo.Values(s.Samplers)
	s.wavetables = lo.Values(s.Wavetables)
}
