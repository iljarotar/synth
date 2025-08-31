package module

import (
	"testing"

	"github.com/iljarotar/synth/calc"
)

func Test_modulate(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		r    calc.Range
		val  float64
		want float64
	}{
		{
			name: "modulating signal is 0",
			x:    1,
			r: calc.Range{
				Min: 0,
				Max: 2,
			},
			val:  0,
			want: 1,
		},
		{
			name: "range is same as output range",
			x:    0,
			r: calc.Range{
				Min: -1,
				Max: 1,
			},
			val:  -0.5,
			want: -0.5,
		},
		{
			name: "range is different than output range",
			x:    0.5,
			r: calc.Range{
				Min: 0,
				Max: 1,
			},
			val:  -0.5,
			want: 0.25,
		},
		{
			name: "modulation exceeds limits of range to lower end",
			x:    400,
			r: calc.Range{
				Min: 0,
				Max: 20000,
			},
			val:  -1,
			want: 0,
		},
		{
			name: "modulation exceeds limits of range to higher end",
			x:    0.75,
			r: calc.Range{
				Min: 0,
				Max: 1,
			},
			val:  1,
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := modulate(tt.x, tt.r, tt.val); got != tt.want {
				t.Errorf("modulate() = %v, want %v", got, tt.want)
			}
		})
	}
}
