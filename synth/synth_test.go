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

func TestSynth_deleteOldModules(t *testing.T) {
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
			name: "delete old and leave others in place",
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
				Envelopes:   module.EnvelopeMap{"env2": {}},
				Filters:     module.FilterMap{"f2": {}},
				Gates:       module.GateMap{"g2": {}},
				Mixers:      module.MixerMap{"m2": {}},
				Noises:      module.NoiseMap{"n2": {}},
				Oscillators: module.OscillatorMap{"o2": {}},
				Pans:        module.PanMap{"p2": {}},
				Samplers:    module.SamplerMap{"s2": {}},
				Sequencers:  module.SequencerMap{"seq2": {}},
				Wavetables:  module.WavetableMap{"w2": {}},
			},
			want: &Synth{
				Out:    "main",
				Volume: 0.5,
				Envelopes: module.EnvelopeMap{
					"env2": env2,
				},
				Filters: module.FilterMap{
					"f2": f2,
				},
				Gates: module.GateMap{
					"g2": g2,
				},
				Mixers: module.MixerMap{
					"m2": m2,
				},
				Noises: module.NoiseMap{
					"n2": n2,
				},
				Oscillators: module.OscillatorMap{
					"o2": o2,
				},
				Pans: module.PanMap{
					"p2": p2,
				},
				Samplers: module.SamplerMap{
					"s2": s2,
				},
				Sequencers: module.SequencerMap{
					"seq2": seq2,
				},
				Wavetables: module.WavetableMap{
					"w2": w2,
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
			tt.s.deleteOldModules(tt.new)
			if diff := cmp.Diff(tt.want, tt.s, cmpopts.IgnoreUnexported(
				Synth{},
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
			); diff != "" {
				t.Errorf("Synth.deleteOldModules diff = %s", diff)
			}
		})
	}
}
