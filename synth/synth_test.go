package synth

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/iljarotar/synth/module"
)

func Test_secondsToStep(t *testing.T) {
	tests := []struct {
		name       string
		seconds    float64
		delta      float64
		sampleRate float64
		want       float64
	}{
		{
			name:       "when seconds is 0 step is delta",
			seconds:    0,
			delta:      1,
			sampleRate: 1,
			want:       1,
		},
		{
			name:       "when seconds is 1 step is delta/sampleRate",
			seconds:    1,
			delta:      0.5,
			sampleRate: 1000,
			want:       0.5 / 1000,
		},
		{
			name:       "when seconds is greater than 1 step is delta/(seconds*sampleRate)",
			seconds:    5,
			delta:      0.1,
			sampleRate: 1000,
			want:       0.1 / 5000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := secondsToStep(tt.seconds, tt.delta, tt.sampleRate); got != tt.want {
				t.Errorf("secondsToStep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSynth_Update(t *testing.T) {
	var (
		env1 = &module.Envelope{}
		env2 = &module.Envelope{}
		f1   = &module.Filter{}
		f2   = &module.Filter{}
		g1   = &module.Gate{}
		g2   = &module.Gate{}
		m1   = &module.Mixer{}
		m2   = &module.Mixer{}
		n1   = &module.Noise{}
		n2   = &module.Noise{}
		o1   = &module.Oscillator{}
		o2   = &module.Oscillator{}
		p1   = &module.Pan{}
		p2   = &module.Pan{}
		s1   = &module.Sampler{}
		s2   = &module.Sampler{}
		seq1 = &module.Sequencer{}
		seq2 = &module.Sequencer{}
		w1   = &module.Wavetable{}
		w2   = &module.Wavetable{}
	)
	tests := []struct {
		name string
		s    *Synth
		new  *Synth
		want *Synth
	}{
		{
			name: "update all modules",
			s: &Synth{
				Out:    "main",
				Volume: 0.5,
				Envelopes: module.EnvelopeMap{
					"env1": env1,
					"env2": env2,
				},
				Filters: module.FilterMap{
					"f1": f1,
					"f2": f2,
				},
				Gates: module.GateMap{
					"g1": g1,
					"g2": g2,
				},
				Mixers: module.MixerMap{
					"m1": m1,
					"m2": m2,
				},
				Noises: module.NoiseMap{
					"n1": n1,
					"n2": n2,
				},
				Oscillators: module.OscillatorMap{
					"o1": o1,
					"o2": o2,
				},
				Pans: module.PanMap{
					"p1": p1,
					"p2": p2,
				},
				Samplers: module.SamplerMap{
					"s1": s1,
					"s2": s2,
				},
				Sequencers: module.SequencerMap{
					"seq1": seq1,
					"seq2": seq2,
				},
				Wavetables: module.WavetableMap{
					"w1": w1,
					"w2": w2,
				},
				Time:         5,
				VolumeMemory: 1,
				sampleRate:   44100,
				volumeStep:   0.1,
				modules: module.ModuleMap{
					"env1": env1,
					"env2": env2,
					"f1":   f1,
					"f2":   f2,
					"g1":   g1,
					"g2":   g2,
					"m1":   m1,
					"m2":   m2,
					"n1":   n1,
					"n2":   n2,
					"o1":   o1,
					"o2":   o2,
					"p1":   p1,
					"p2":   p2,
					"s1":   s1,
					"s2":   s2,
					"seq1": seq1,
					"seq2": seq2,
					"w1":   w1,
					"w2":   w2,
				},
				envelopes:   []*module.Envelope{env1, env2},
				filters:     []*module.Filter{f1, f2},
				gates:       []*module.Gate{g1, g2},
				mixers:      []*module.Mixer{m1, m2},
				noises:      []*module.Noise{n1, n2},
				oscillators: []*module.Oscillator{o1, o2},
				pans:        []*module.Pan{p1, p2},
				samplers:    []*module.Sampler{s1, s2},
				sequencers:  []*module.Sequencer{seq1, seq2},
				wavetables:  []*module.Wavetable{w1, w2},
			},
			new: &Synth{
				Envelopes: module.EnvelopeMap{
					"env2": {
						Gate:    "new-gate",
						Attack:  1,
						Decay:   1,
						Release: 1,
						Peak:    100,
						Level:   1,
					},
				},
				Filters: module.FilterMap{
					"f2": {
						In:    "new-in",
						Type:  "BandPass",
						Freq:  200,
						Width: 10,
						CV:    "new-cv",
						Mod:   "new-mod",
					},
				},
				Gates: module.GateMap{
					"g2": {
						CV:     "new-cv",
						BPM:    20,
						Mod:    "new-mod",
						Signal: []float64{1},
					},
				},
				Mixers: module.MixerMap{
					"m2": {
						CV:   "new-mod",
						Gain: 1,
						Mod:  "new-mod",
						In: map[string]float64{
							"new-in": 1,
						},
					},
				},
				Noises: module.NoiseMap{"n2": {}},
				Oscillators: module.OscillatorMap{
					"o2": {
						Module: module.Module{},
						Type:   "Sine",
						Freq:   200,
						CV:     "new-cv",
						Mod:    "new-mod",
						Phase:  0.5,
					},
				},
				Pans: module.PanMap{
					"p2": {
						Pan: 1,
						Mod: "new-mod",
						In:  "new-in",
					},
				},
				Samplers: module.SamplerMap{
					"s2": {
						In:      "new-in",
						Trigger: "new-trigger",
					},
				},
				Sequencers: module.SequencerMap{
					"seq2": {
						Sequence:  []string{"a_4"},
						Trigger:   "new-trigger",
						Pitch:     440,
						Transpose: 1,
						Randomize: true,
					},
				},
				Wavetables: module.WavetableMap{
					"w2": {
						Freq:   300,
						CV:     "new-cv",
						Mod:    "new-mod",
						Signal: []float64{2},
					},
				},
			},
			want: &Synth{
				Out:    "main",
				Volume: 0.5,
				Envelopes: module.EnvelopeMap{
					"env2": {
						Gate: "new-gate",
					},
				},
				Filters: module.FilterMap{
					"f2": {
						In:   "new-in",
						Type: "BandPass",
						CV:   "new-cv",
						Mod:  "new-mod",
					},
				},
				Gates: module.GateMap{
					"g2": {
						CV:     "new-cv",
						Mod:    "new-mod",
						Signal: []float64{1},
					},
				},
				Mixers: module.MixerMap{
					"m2": {
						CV:  "new-mod",
						Mod: "new-mod",
						In: map[string]float64{
							"new-in": 0,
						},
					},
				},
				Noises: module.NoiseMap{
					"n2": n2,
				},
				Oscillators: module.OscillatorMap{
					"o2": {
						Module: module.Module{},
						Type:   "Sine",
						CV:     "new-cv",
						Mod:    "new-mod",
					},
				},
				Pans: module.PanMap{
					"p2": {
						Mod: "new-mod",
						In:  "new-in",
					},
				},
				Samplers: module.SamplerMap{
					"s2": {
						In:      "new-in",
						Trigger: "new-trigger",
					},
				},
				Sequencers: module.SequencerMap{
					"seq2": {
						Sequence:  []string{"a_4"},
						Trigger:   "new-trigger",
						Pitch:     440,
						Transpose: 1,
						Randomize: true,
					},
				},
				Wavetables: module.WavetableMap{
					"w2": {
						CV:     "new-cv",
						Mod:    "new-mod",
						Signal: []float64{2},
					},
				},
				Time:         5,
				VolumeMemory: 1,
				sampleRate:   44100,
				volumeStep:   0.1,
				modules: module.ModuleMap{
					"env2": env2,
					"f2":   f2,
					"g2":   g2,
					"m2":   m2,
					"n2":   n2,
					"o2":   o2,
					"p2":   p2,
					"s2":   s2,
					"seq2": seq2,
					"w2":   w2,
				},
				envelopes:   []*module.Envelope{env2},
				filters:     []*module.Filter{f2},
				gates:       []*module.Gate{g2},
				mixers:      []*module.Mixer{m2},
				noises:      []*module.Noise{n2},
				oscillators: []*module.Oscillator{o2},
				pans:        []*module.Pan{p2},
				samplers:    []*module.Sampler{s2},
				sequencers:  []*module.Sequencer{seq2},
				wavetables:  []*module.Wavetable{w2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Update(tt.new)
			if diff := cmp.Diff(tt.want, tt.s,
				cmpopts.IgnoreUnexported(
					module.Module{},
					module.Envelope{},
					module.Filter{},
					module.Gate{},
					module.Mixer{},
					module.Noise{},
					module.Oscillator{},
					module.Pan{},
					module.Sampler{},
					module.Sequencer{},
					module.Wavetable{},
				),
				cmp.AllowUnexported(Synth{}),
			); diff != "" {
				t.Errorf("Synth.Update() diff = %s", diff)
			}
		})
	}
}

func TestSynth_initializeEmptyMaps(t *testing.T) {
	tests := []struct {
		name string
		s    *Synth
		want *Synth
	}{
		{
			name: "initialize empty",
			s:    &Synth{},
			want: &Synth{
				Envelopes:   module.EnvelopeMap{},
				Filters:     module.FilterMap{},
				Gates:       module.GateMap{},
				Mixers:      module.MixerMap{},
				Noises:      module.NoiseMap{},
				Oscillators: module.OscillatorMap{},
				Pans:        module.PanMap{},
				Samplers:    module.SamplerMap{},
				Sequencers:  module.SequencerMap{},
				Wavetables:  module.WavetableMap{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.initializeEmptyMaps()
			if diff := cmp.Diff(tt.want, tt.s, cmp.AllowUnexported(Synth{})); diff != "" {
				t.Errorf("Synth.initializeEmptyMaps() diff = %s", diff)
			}
		})
	}
}
