package module

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_normalizeWavetable(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   []float64
	}{
		{
			name:   "within range [-1,1] no change",
			values: []float64{-1, -0.3, 0, 0.2, 1},
			want:   []float64{-1, -0.3, 0, 0.2, 1},
		},
		{
			name:   "positive values exceeding 1",
			values: []float64{0, 1, 2, 3, 4},
			want:   []float64{0, 0.25, 0.5, 0.75, 1},
		},
		{
			name:   "negative values lower than -1",
			values: []float64{0, -1, -2, -3, -4},
			want:   []float64{0, -0.25, -0.5, -0.75, -1},
		},
		{
			name:   "asymmetrically spread values",
			values: []float64{-1, 0, 1, 2, 3, 4},
			want:   []float64{-0.25, 0, 0.25, 0.5, 0.75, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeWavetable(tt.values)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("normalizeWavetable() diff %v", diff)
			}
		})
	}
}
