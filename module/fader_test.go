package module

import (
	"testing"
)

func Test_fader_initialize(t *testing.T) {
	tests := []struct {
		name       string
		duration   float64
		sampleRate float64
		f          *fader
		wantStep   float64
	}{
		{
			name:       "duration is zero",
			duration:   0,
			sampleRate: 44100,
			f: &fader{
				current: 440,
				target:  330,
			},
			wantStep: -110,
		},
		{
			name:       "sample rate is zero",
			duration:   1,
			sampleRate: 0,
			f: &fader{
				current: 440,
				target:  330,
			},
			wantStep: -110,
		},
		{
			name:       "both non-zero",
			duration:   5,
			sampleRate: 2000,
			f: &fader{
				current: 400,
				target:  300,
			},
			wantStep: -0.01,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.initialize(tt.duration, tt.sampleRate)
			if tt.f.step != tt.wantStep {
				t.Errorf("fader.initialize() step = %v, want %v", tt.f.step, tt.wantStep)
			}
		})
	}
}

func Test_fader_fade(t *testing.T) {
	tests := []struct {
		name string
		f    *fader
		want float64
	}{
		{
			name: "target equals current",
			f: &fader{
				current: 400,
				target:  400,
				step:    10,
			},
			want: 400,
		},
		{
			name: "current lower than target",
			f: &fader{
				current: 400,
				target:  500,
				step:    10,
			},
			want: 410,
		},
		{
			name: "current higher than target",
			f: &fader{
				current: 400,
				target:  300,
				step:    -10,
			},
			want: 390,
		},
		{
			name: "current lower than target but exceeds it after update",
			f: &fader{
				current: 395,
				target:  400,
				step:    10,
			},
			want: 400,
		},
		{
			name: "current higher than target but lower after update",
			f: &fader{
				current: 305,
				target:  300,
				step:    -10,
			},
			want: 300,
		},
		{
			name: "current higher than target but step is positive",
			f: &fader{
				current: 305,
				target:  300,
				step:    10,
			},
			want: 300,
		},
		{
			name: "current lower than target but step is negative",
			f: &fader{
				current: 295,
				target:  300,
				step:    -10,
			},
			want: 300,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.fade(); got != tt.want {
				t.Errorf("fader.fade() = %v, want %v", got, tt.want)
			}
		})
	}
}
