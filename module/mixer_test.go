package module

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMixer_Step(t *testing.T) {
	tests := []struct {
		name         string
		m            *Mixer
		modules      ModuleMap
		want         Output
		wantIntegral float64
	}{
		{
			name: "no input",
			m: &Mixer{
				Gain:       0,
				In:         map[string]float64{},
				sampleRate: 1,
			},
			modules: ModuleMap{
				"in": &Module{},
			},
			want:         Output{},
			wantIntegral: 0,
		},
		{
			name: "no module found",
			m: &Mixer{
				Gain: 1,
				In: map[string]float64{
					"sine": 1,
				},
				sampleRate: 1,
			},
			modules: ModuleMap{
				"in": &Oscillator{
					Module: Module{},
				},
			},
			want:         Output{},
			wantIntegral: 0,
		},
		{
			name: "input gain 0",
			m: &Mixer{
				Gain: 1,
				In: map[string]float64{
					"in": 0,
				},
				sampleRate: 1,
			},
			modules: ModuleMap{
				"in": &Module{
					current: Output{
						Mono:  1,
						Left:  0.5,
						Right: 0.5,
					},
				},
			},
			want:         Output{},
			wantIntegral: 0,
		},
		{
			name: "input",
			m: &Mixer{
				Gain: 1,
				In: map[string]float64{
					"in": 1,
				},
				sampleRate: 1,
			},
			modules: ModuleMap{
				"in": &Module{
					current: Output{
						Mono:  1,
						Left:  0.5,
						Right: 0.5,
					},
				},
			},
			want: Output{
				Mono:  1,
				Left:  0.5,
				Right: 0.5,
			},
			wantIntegral: 0.5,
		},
		{
			name: "mod",
			m: &Mixer{
				Gain: 0.5,
				Mod:  "lfo",
				In: map[string]float64{
					"in": 1,
				},
				sampleRate: 1,
			},
			modules: ModuleMap{
				"in": &Module{
					current: Output{
						Mono:  1,
						Left:  0.5,
						Right: 0.5,
					},
				},
				"lfo": &Module{
					current: Output{
						Mono:  0.5,
						Left:  0.25,
						Right: 0.25,
					},
				},
			},
			want: Output{
				Mono:  1 * 0.75,
				Left:  0.5 * 0.75,
				Right: 0.5 * 0.75,
			},
			wantIntegral: 0.5 * 0.75,
		},
		{
			name: "cv",
			m: &Mixer{
				Gain: 0.5,
				Mod:  "in",
				CV:   "lfo",
				In: map[string]float64{
					"in": 1,
				},
				sampleRate: 1,
			},
			modules: ModuleMap{
				"in": &Module{
					current: Output{
						Mono:  1,
						Left:  0.5,
						Right: 0.5,
					},
				},
				"lfo": &Module{
					current: Output{
						Mono:  0.5,
						Left:  0.25,
						Right: 0.25,
					},
				},
			},
			want: Output{
				Mono:  1 * 0.75,
				Left:  0.5 * 0.75,
				Right: 0.5 * 0.75,
			},
			wantIntegral: 0.5 * 0.75,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Step(tt.modules)

			if diff := cmp.Diff(tt.want, tt.m.Current()); diff != "" {
				t.Errorf("Mixer.Step() diff = %s", diff)
			}
		})
	}
}
