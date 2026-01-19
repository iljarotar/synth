package synth

import (
	"slices"

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

	Envelopes   module.EnvelopeMap   `yaml:"envelopes"`
	Filters     module.FilterMap     `yaml:"filters"`
	Gates       module.GateMap       `yaml:"gates"`
	Mixers      module.MixerMap      `yaml:"mixers"`
	Noises      module.NoiseMap      `yaml:"noises"`
	Oscillators module.OscillatorMap `yaml:"oscillators"`
	Pans        module.PanMap        `yaml:"pans"`
	Samplers    module.SamplerMap    `yaml:"samplers"`
	Sequencers  module.SequencerMap  `yaml:"sequencers"`
	Wavetables  module.WavetableMap  `yaml:"wavetables"`

	Time              float64
	VolumeMemory      float64
	sampleRate        float64
	volumeStep        float64
	notifyFadeoutChan chan<- bool
	logger            *log.Logger
	modules           module.ModuleMap

	envelopes   []*module.Envelope
	filters     []*module.Filter
	gates       []*module.Gate
	mixers      []*module.Mixer
	noises      []*module.Noise
	oscillators []*module.Oscillator
	pans        []*module.Pan
	samplers    []*module.Sampler
	sequencers  []*module.Sequencer
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
	s.initializeEmptyMaps()
	s.makeModulesMap()
	s.flattenModules()

	if err := s.Filters.Initialize(sampleRate); err != nil {
		return err
	}
	if err := s.Mixers.Initialize(sampleRate); err != nil {
		return err
	}
	if err := s.Oscillators.Initialize(sampleRate); err != nil {
		return err
	}
	if err := s.Sequencers.Initialize(); err != nil {
		return err
	}

	s.Envelopes.Initialize(sampleRate)
	s.Gates.Initialze(sampleRate)
	s.Pans.Initialize(sampleRate)
	s.Wavetables.Initialize(sampleRate)

	return nil
}

func (s *Synth) Update(from *Synth) error {
	if from == nil {
		return nil
	}

	s.deleteOldModules(from)
	s.addNewModules(from)
	s.updateModules(from)

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
	for _, e := range s.envelopes {
		if e == nil {
			continue
		}
		e.Step(s.Time, s.modules)
	}
	for _, f := range s.filters {
		if f == nil {
			continue
		}
		f.Step(s.modules)
	}
	for _, g := range s.gates {
		if g == nil {
			continue
		}
		g.Step(s.modules)
	}
	for _, m := range s.mixers {
		if m == nil {
			continue
		}
		m.Step(s.modules)
	}
	for _, n := range s.noises {
		if n == nil {
			continue
		}
		n.Step()
	}
	for _, osc := range s.oscillators {
		if osc == nil {
			continue
		}
		osc.Step(s.modules)
	}
	for _, p := range s.pans {
		if p == nil {
			continue
		}
		p.Step(s.modules)
	}
	for _, smplr := range s.samplers {
		if smplr == nil {
			continue
		}
		smplr.Step(s.modules)
	}
	for _, seq := range s.sequencers {
		if seq == nil {
			continue
		}
		seq.Step(s.modules)
	}
	for _, w := range s.wavetables {
		if w == nil {
			continue
		}
		w.Step(s.modules)
	}

	s.Time += 1 / s.sampleRate
}

func secondsToStep(seconds, delta, sampleRate float64) float64 {
	if sampleRate == 0 || seconds == 0 {
		return delta
	}
	return delta / (seconds * sampleRate)
}

func (s *Synth) makeModulesMap() {
	s.modules = module.ModuleMap{}

	for name, e := range s.Envelopes {
		if e == nil {
			continue
		}
		s.modules[name] = e
	}
	for name, f := range s.Filters {
		if f == nil {
			continue
		}
		s.modules[name] = f
	}
	for name, g := range s.Gates {
		if g == nil {
			continue
		}
		s.modules[name] = g
	}
	for name, m := range s.Mixers {
		if m == nil {
			continue
		}
		s.modules[name] = m
	}
	for name, n := range s.Noises {
		if n == nil {
			continue
		}
		s.modules[name] = n
	}
	for name, osc := range s.Oscillators {
		if osc == nil {
			continue
		}
		s.modules[name] = osc
	}
	for name, p := range s.Pans {
		if p == nil {
			continue
		}
		s.modules[name] = p
	}
	for name, smplr := range s.Samplers {
		if smplr == nil {
			continue
		}
		s.modules[name] = smplr
	}
	for name, seq := range s.Sequencers {
		if seq == nil {
			continue
		}
		s.modules[name] = seq
	}
	for name, w := range s.Wavetables {
		if w == nil {
			continue
		}
		s.modules[name] = w
	}
}

func (s *Synth) flattenModules() {
	s.envelopes = lo.Values(s.Envelopes)
	s.filters = lo.Values(s.Filters)
	s.gates = lo.Values(s.Gates)
	s.mixers = lo.Values(s.Mixers)
	s.noises = lo.Values(s.Noises)
	s.oscillators = lo.Values(s.Oscillators)
	s.pans = lo.Values(s.Pans)
	s.samplers = lo.Values(s.Samplers)
	s.sequencers = lo.Values(s.Sequencers)
	s.wavetables = lo.Values(s.Wavetables)
}

func (s *Synth) deleteOldModules(new *Synth) {
	for name, env := range s.Envelopes {
		if _, ok := new.Envelopes[name]; !ok {
			delete(s.Envelopes, name)
			delete(s.modules, name)
			s.envelopes = slices.DeleteFunc(s.envelopes, func(e *module.Envelope) bool {
				return env == e
			})
		}
	}
	for name, filter := range s.Filters {
		if _, ok := new.Filters[name]; !ok {
			delete(s.Filters, name)
			delete(s.modules, name)
			s.filters = slices.DeleteFunc(s.filters, func(f *module.Filter) bool {
				return filter == f
			})
		}
	}
	for name, gate := range s.Gates {
		if _, ok := new.Gates[name]; !ok {
			delete(s.Gates, name)
			delete(s.modules, name)
			s.gates = slices.DeleteFunc(s.gates, func(g *module.Gate) bool {
				return gate == g
			})
		}
	}
	for name, mixer := range s.Mixers {
		if _, ok := new.Mixers[name]; !ok {
			delete(s.Mixers, name)
			delete(s.modules, name)
			s.mixers = slices.DeleteFunc(s.mixers, func(m *module.Mixer) bool {
				return mixer == m
			})
		}
	}
	for name, noise := range s.Noises {
		if _, ok := new.Noises[name]; !ok {
			delete(s.Noises, name)
			delete(s.modules, name)
			s.noises = slices.DeleteFunc(s.noises, func(n *module.Noise) bool {
				return noise == n
			})
		}
	}
	for name, osc := range s.Oscillators {
		if _, ok := new.Oscillators[name]; !ok {
			delete(s.Oscillators, name)
			delete(s.modules, name)
			s.oscillators = slices.DeleteFunc(s.oscillators, func(o *module.Oscillator) bool {
				return osc == o
			})
		}
	}
	for name, pan := range s.Pans {
		if _, ok := new.Pans[name]; !ok {
			delete(s.Pans, name)
			delete(s.modules, name)
			s.pans = slices.DeleteFunc(s.pans, func(p *module.Pan) bool {
				return pan == p
			})
		}
	}
	for name, sampler := range s.Samplers {
		if _, ok := new.Samplers[name]; !ok {
			delete(s.Samplers, name)
			delete(s.modules, name)
			s.samplers = slices.DeleteFunc(s.samplers, func(smplr *module.Sampler) bool {
				return sampler == smplr
			})
		}
	}
	for name, seq := range s.Sequencers {
		if _, ok := new.Sequencers[name]; !ok {
			delete(s.Sequencers, name)
			delete(s.modules, name)
			s.sequencers = slices.DeleteFunc(s.sequencers, func(sq *module.Sequencer) bool {
				return seq == sq
			})
		}
	}
	for name, wt := range s.Wavetables {
		if _, ok := new.Wavetables[name]; !ok {
			delete(s.Wavetables, name)
			delete(s.modules, name)
			s.wavetables = slices.DeleteFunc(s.wavetables, func(w *module.Wavetable) bool {
				return wt == w
			})
		}
	}
}

func (s *Synth) addNewModules(new *Synth) {
	for name, env := range new.Envelopes {
		if _, ok := s.Envelopes[name]; !ok {
			s.Envelopes[name] = env
			s.envelopes = append(s.envelopes, env)
			s.modules[name] = env
		}
	}
	for name, filter := range new.Filters {
		if _, ok := s.Filters[name]; !ok {
			s.Filters[name] = filter
			s.filters = append(s.filters, filter)
			s.modules[name] = filter
		}
	}
	for name, gate := range new.Gates {
		if _, ok := s.Gates[name]; !ok {
			s.Gates[name] = gate
			s.gates = append(s.gates, gate)
			s.modules[name] = gate
		}
	}
	for name, mixer := range new.Mixers {
		if _, ok := s.Mixers[name]; !ok {
			s.Mixers[name] = mixer
			s.mixers = append(s.mixers, mixer)
			s.modules[name] = mixer
		}
	}
	for name, noise := range new.Noises {
		if _, ok := s.Noises[name]; !ok {
			s.Noises[name] = noise
			s.noises = append(s.noises, noise)
			s.modules[name] = noise
		}
	}
	for name, osc := range new.Oscillators {
		if _, ok := s.Oscillators[name]; !ok {
			s.Oscillators[name] = osc
			s.oscillators = append(s.oscillators, osc)
			s.modules[name] = osc
		}
	}
	for name, pan := range new.Pans {
		if _, ok := s.Pans[name]; !ok {
			s.Pans[name] = pan
			s.pans = append(s.pans, pan)
			s.modules[name] = pan
		}
	}
	for name, sampler := range new.Samplers {
		if _, ok := s.Samplers[name]; !ok {
			s.Samplers[name] = sampler
			s.samplers = append(s.samplers, sampler)
			s.modules[name] = sampler
		}
	}
	for name, seq := range new.Sequencers {
		if _, ok := s.Sequencers[name]; !ok {
			s.Sequencers[name] = seq
			s.sequencers = append(s.sequencers, seq)
			s.modules[name] = seq
		}
	}
	for name, wt := range new.Wavetables {
		if _, ok := s.Wavetables[name]; !ok {
			s.Wavetables[name] = wt
			s.wavetables = append(s.wavetables, wt)
			s.modules[name] = wt
		}
	}
}

func (s *Synth) updateModules(new *Synth) {
	for name, env := range s.Envelopes {
		if newEnv, ok := new.Envelopes[name]; ok {
			env.Update(newEnv)
		}
	}
	for name, filter := range s.Filters {
		if newFilter, ok := new.Filters[name]; ok {
			filter.Update(newFilter)
		}
	}
	for name, gate := range s.Gates {
		if newGate, ok := new.Gates[name]; ok {
			gate.Update(newGate)
		}
	}
	for name, mixer := range s.Mixers {
		if newMixer, ok := new.Mixers[name]; ok {
			mixer.Update(newMixer)
		}
	}
	for name, osc := range s.Oscillators {
		if newOsc, ok := new.Oscillators[name]; ok {
			osc.Update(newOsc)
		}
	}
	for name, pan := range s.Pans {
		if newPan, ok := new.Pans[name]; ok {
			pan.Update(newPan)
		}
	}
	for name, sampler := range s.Samplers {
		if newSampler, ok := new.Samplers[name]; ok {
			sampler.Update(newSampler)
		}
	}
	for name, seq := range s.Sequencers {
		if newSeq, ok := new.Sequencers[name]; ok {
			seq.Update(newSeq)
		}
	}
	for name, wt := range s.Wavetables {
		if newWt, ok := new.Wavetables[name]; ok {
			wt.Update(newWt)
		}
	}
}

func (s *Synth) initializeEmptyMaps() {
	if s.Envelopes == nil {
		s.Envelopes = module.EnvelopeMap{}
	}
	if s.Filters == nil {
		s.Filters = module.FilterMap{}
	}
	if s.Gates == nil {
		s.Gates = module.GateMap{}
	}
	if s.Mixers == nil {
		s.Mixers = module.MixerMap{}
	}
	if s.Noises == nil {
		s.Noises = module.NoiseMap{}
	}
	if s.Oscillators == nil {
		s.Oscillators = module.OscillatorMap{}
	}
	if s.Pans == nil {
		s.Pans = module.PanMap{}
	}
	if s.Samplers == nil {
		s.Samplers = module.SamplerMap{}
	}
	if s.Sequencers == nil {
		s.Sequencers = module.SequencerMap{}
	}
	if s.Wavetables == nil {
		s.Wavetables = module.WavetableMap{}
	}
}
