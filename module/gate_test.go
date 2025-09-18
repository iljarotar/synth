package module

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/iljarotar/synth/calc"
)

func TestGate_initialze(t *testing.T) {
	tests := []struct {
		name       string
		g          *Gate
		wantSignal []float64
	}{
		{
			name: "convert to 0s and 1s only",
			g: &Gate{
				Signal: []float64{0, 1, -400.1, 2},
			},
			wantSignal: []float64{-1, 1, -1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.initialze(44100)

			if diff := cmp.Diff(tt.wantSignal, tt.g.Signal); diff != "" {
				t.Errorf("Gate.initialize() signal diff = %s", diff)
			}
		})
	}
}

func TestGate_Step(t *testing.T) {
	sampleRate := 44100.0

	tests := []struct {
		name    string
		g       *Gate
		modules ModuleMap
		want    float64
		wantIdx float64
	}{
		{
			name: "bpm 0",
			g: &Gate{
				Module:     Module{},
				BPM:        0,
				Signal:     []float64{1, -1},
				sampleRate: sampleRate,
				idx:        0,
			},
			want:    1,
			wantIdx: 0,
		},
		{
			name: "no mod, no cv",
			g: &Gate{
				BPM:        60,
				Signal:     []float64{-1, 1, -1, 1},
				sampleRate: sampleRate,
				idx:        0,
			},
			want:    -1,
			wantIdx: 1 / sampleRate,
		},
		{
			name: "cv",
			g: &Gate{
				BPM:        60,
				Signal:     []float64{-1, 1, -1, 1},
				sampleRate: sampleRate,
				idx:        1,
				CV:         "cv",
			},
			modules: ModuleMap{
				"cv": &Module{
					current: Output{
						Mono: calc.Transpose(120, bpmLimits, outputLimits),
					},
				},
			},
			want:    1,
			wantIdx: 1 + 2/sampleRate,
		},
		{
			name: "mod",
			g: &Gate{
				BPM:        60,
				Signal:     []float64{-1, 1, -1, 1},
				sampleRate: sampleRate,
				idx:        2,
				Mod:        "mod",
			},
			modules: ModuleMap{
				"mod": &Module{
					current: Output{
						Mono: -0.03,
					},
				},
			},
			want:    -1,
			wantIdx: 2 + 0.5/sampleRate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.Step(tt.modules)

			if tt.g.current.Mono != tt.want {
				t.Errorf("Gate.Step() = %v, want %v", tt.g.current.Mono, tt.want)
			}
			if tt.g.idx != tt.wantIdx {
				t.Errorf("Gate.Step() idx = %v, wantIdx %v", tt.g.idx, tt.wantIdx)
			}
		})
	}
}

func Test_samplesPerBeat(t *testing.T) {
	tests := []struct {
		name       string
		sampleRate float64
		bpm        float64
		want       float64
	}{
		{
			name:       "bpm 0",
			sampleRate: 0,
			bpm:        0,
			want:       math.Inf(1),
		},
		{
			name:       "bpm 60",
			sampleRate: 44100,
			bpm:        60,
			want:       44100,
		},
		{
			name:       "bpm 120",
			sampleRate: 44100,
			bpm:        120,
			want:       22050,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := samplesPerBeat(tt.sampleRate, tt.bpm); got != tt.want {
				t.Errorf("samplesPerBeat() = %v, want %v", got, tt.want)
			}
		})
	}
}
