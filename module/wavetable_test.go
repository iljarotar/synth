package module

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWavetable_initialze(t *testing.T) {
	tests := []struct {
		name string
		w    *Wavetable
		want []float64
	}{
		{
			name: "limit exceeding values",
			w: &Wavetable{
				Signal: []float64{-0.75, 0, -2, 0.75, 2},
			},
			want: []float64{-0.75, 0, -1, 0.75, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.initialze(0)

			if diff := cmp.Diff(tt.want, tt.w.Signal); diff != "" {
				t.Errorf("Wavetable.initialize() diff = %s", diff)
			}
		})
	}
}

func TestWavetable_Step(t *testing.T) {
	sampleRate := 44100.0

	tests := []struct {
		name    string
		w       *Wavetable
		modules ModuleMap
		want    float64
		wantIdx float64
	}{
		{
			name: "no mod no cv",
			w: &Wavetable{
				Freq:       2,
				Signal:     []float64{1, 0, -1, 0},
				sampleRate: sampleRate,
				idx:        44100.0 / 8,
			},
			modules: ModuleMap{},
			want:    0,
			wantIdx: 44100.0/8 + 8/44100.0,
		},
		{
			name: "cv",
			w: &Wavetable{
				Freq:       2,
				Signal:     []float64{1, 0, -1, 0},
				CV:         "cv",
				sampleRate: sampleRate,
				idx:        0,
			},
			modules: ModuleMap{
				"cv": &Module{
					current: Output{
						Mono: 0,
					},
				},
			},
			want:    1,
			wantIdx: 4 * freqLimits.Max / (2 * sampleRate),
		},
		{
			name: "mod",
			w: &Wavetable{
				Freq:       2,
				Signal:     []float64{1, 0, -1, 0},
				Mod:        "mod",
				sampleRate: sampleRate,
				idx:        2.5,
			},
			modules: ModuleMap{
				"mod": &Module{
					current: Output{
						Mono: 1,
					},
				},
			},
			want:    -1,
			wantIdx: 2.5 + 16/sampleRate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.Step(tt.modules)
		})
	}
}
