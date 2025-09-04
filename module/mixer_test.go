package module

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMixer_Step(t *testing.T) {
	tests := []struct {
		name         string
		m            *Mixer
		modules      ModulesMap
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
			modules: ModulesMap{
				"sine": &Oscillator{
					Module: Module{},
				},
			},
			want:         Output{},
			wantIntegral: 0,
		},
		{
			name: "no module found",
			m: &Mixer{
				Gain: 1,
				In: map[string]float64{
					"square": 1,
				},
				sampleRate: 1,
			},
			modules: ModulesMap{
				"sine": &Oscillator{
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
					"sine": 0,
				},
				sampleRate: 1,
			},
			modules: ModulesMap{
				"sine": &Oscillator{
					Module: Module{
						current: Output{
							Mono:  1,
							Left:  0.5,
							Right: 0.5,
						},
					},
				},
			},
			want:         Output{},
			wantIntegral: 0,
		},
		{
			name: "oscillator input",
			m: &Mixer{
				Gain: 1,
				In: map[string]float64{
					"sine": 1,
				},
				sampleRate: 1,
			},
			modules: ModulesMap{
				"sine": &Oscillator{
					Module: Module{
						current: Output{
							Mono:  1,
							Left:  0.5,
							Right: 0.5,
						},
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
			name: "mixer input",
			m: &Mixer{
				Gain: 1,
				In: map[string]float64{
					"mix": 1,
				},
				sampleRate: 1,
			},
			modules: ModulesMap{
				"mix": &Mixer{
					Module: Module{
						current: Output{
							Mono:  1,
							Left:  0.5,
							Right: 0.5,
						},
					},
					Gain: 1,
				},
			},
			want: Output{
				Mono:  1,
				Left:  0.5,
				Right: 0.5,
			},
			wantIntegral: 0.5,
		},
		// TODO: add cv and modulation tests
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Step(tt.modules)
		})
		if diff := cmp.Diff(tt.want, tt.m.Current()); diff != "" {
			t.Errorf("Mixer.Step() diff = %s", diff)
		}
		if tt.wantIntegral != tt.m.Integral() {
			t.Errorf("Mixer.Step() integral = %v, want %v", tt.m.Integral(), tt.wantIntegral)
		}
	}
}
