package module

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_comb_initialize(t *testing.T) {
	tests := []struct {
		name       string
		time       float64
		sampleRate float64
		wantY      []float64
	}{
		{
			name:       "time zero",
			time:       0,
			sampleRate: 400,
			wantY:      []float64{},
		},
		{
			name:       "time one second",
			time:       1,
			sampleRate: 400,
			wantY:      make([]float64, 400),
		},
		{
			name:       "time 0.5 seconds",
			time:       0.5,
			sampleRate: 400,
			wantY:      make([]float64, 200),
		},
		{
			name:       "ceil fractional",
			time:       0.3,
			sampleRate: 454,
			wantY:      make([]float64, 137),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &comb{}
			c.initialize(tt.time, tt.sampleRate)
			if diff := cmp.Diff(tt.wantY, c.y); diff != "" {
				t.Errorf("comb.initialize() diff y = %s", diff)
			}
		})
	}
}

func Test_comb_step(t *testing.T) {
	tests := []struct {
		name    string
		c       *comb
		x       float64
		gain    float64
		want    float64
		wantY   []float64
		wantIdx int
	}{
		{
			name: "empty comb",
			c: &comb{
				y: []float64{},
			},
			x:     1,
			gain:  0.5,
			want:  0.5,
			wantY: []float64{},
		},
		{
			name: "first step",
			c: &comb{
				y: []float64{0, 0, 0, 0, 0},
			},
			x:       1,
			gain:    0.25,
			want:    0.75,
			wantY:   []float64{0.75, 0, 0, 0, 0},
			wantIdx: 1,
		},
		{
			name: "shift idx",
			c: &comb{
				y:   []float64{0.5, 0, 0.25, 0, 0},
				idx: 2,
			},
			x:       1,
			gain:    0.25,
			want:    0.8125,
			wantY:   []float64{0.5, 0, 0.8125, 0, 0},
			wantIdx: 3,
		},
		{
			name: "idx should not exceed len(y)",
			c: &comb{
				y:   []float64{0.5, 0, 0.25, 0, 0.5},
				idx: 4,
			},
			x:       1,
			gain:    0.25,
			want:    0.875,
			wantY:   []float64{0.5, 0, 0.25, 0, 0.875},
			wantIdx: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.step(tt.x, tt.gain); got != tt.want {
				t.Errorf("comb.step() = %v, want %v", got, tt.want)
			}
			if tt.c.idx != tt.wantIdx {
				t.Errorf("comb.step() idx = %v, want %v", tt.c.idx, tt.wantIdx)
			}
			if diff := cmp.Diff(tt.wantY, tt.c.y); diff != "" {
				t.Errorf("comb.step() diff y = %s", diff)
			}
		})
	}
}

func Test_comb_update(t *testing.T) {
	tests := []struct {
		name    string
		c       *comb
		time    float64
		wantY   []float64
		wantIdx int
	}{
		{
			name: "empty comb",
			c: &comb{
				sampleRate: 400,
			},
			time:  1,
			wantY: make([]float64, 400),
		},
		{
			name: "longer time than before",
			c: &comb{
				y:          []float64{0.5, 0, 0, 0, 0.5},
				sampleRate: 6,
				idx:        2,
			},
			time:    1,
			wantY:   []float64{0.5, 0, 0, 0, 0.5, 0},
			wantIdx: 2,
		},
		{
			name: "shorter time than before",
			c: &comb{
				y:          []float64{0.5, 0, 0, 0, 0.5},
				sampleRate: 6,
				idx:        4,
			},
			time:    0.5,
			wantY:   []float64{0.5, 0, 0},
			wantIdx: 2,
		},
		{
			name: "update to zero time",
			c: &comb{
				y:          []float64{0.5, 0, 0, 0, 0.5},
				sampleRate: 6,
				idx:        4,
			},
			time:    0,
			wantY:   []float64{},
			wantIdx: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.update(tt.time)
			if diff := cmp.Diff(tt.wantY, tt.c.y); diff != "" {
				t.Errorf("comb.update() diff y = %s", diff)
			}
			if tt.c.idx != tt.wantIdx {
				t.Errorf("comb.update() idx = %v, want %v", tt.c.idx, tt.wantIdx)
			}
		})
	}
}
