package calc

import (
	"testing"
)

func TestPercentage(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		r    Range
		want float64
	}{
		{
			name: "minimum is 0",
			x:    100,
			r: Range{
				Min: 100,
				Max: 200,
			},
			want: 0,
		},
		{
			name: "middle is 0.5",
			x:    150,
			r: Range{
				Min: 100,
				Max: 200,
			},
			want: 0.5,
		},
		{
			name: "maximum is 1",
			x:    200,
			r: Range{
				Min: 100,
				Max: 200,
			},
			want: 1,
		},
		{
			name: "higher than max",
			x:    300,
			r: Range{
				Min: 100,
				Max: 200,
			},
			want: 2,
		},
		{
			name: "lower than min",
			x:    -100,
			r: Range{
				Min: 100,
				Max: 200,
			},
			want: -2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Percentage(tt.x, tt.r); got != tt.want {
				t.Errorf("Percentage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLimit(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		r    Range
		want float64
	}{
		{
			name: "idempotent on the lower end",
			x:    100,
			r: Range{
				Min: 100,
				Max: 200,
			},
			want: 100,
		},
		{
			name: "idempotent on the upper end",
			x:    200,
			r: Range{
				Min: 100,
				Max: 200,
			},
			want: 200,
		},
		{
			name: "upper limit",
			x:    300,
			r: Range{
				Min: 100,
				Max: 200,
			},
			want: 200,
		},
		{
			name: "lower limit",
			x:    0,
			r: Range{
				Min: 100,
				Max: 200,
			},
			want: 100,
		},
		{
			name: "don't limit",
			x:    0.1,
			r: Range{
				Min: 0,
				Max: 1,
			},
			want: 0.1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Limit(tt.x, tt.r); got != tt.want {
				t.Errorf("Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranspose(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		from Range
		to   Range
		want float64
	}{
		{
			name: "equal ranges",
			x:    6,
			from: Range{
				Min: 1,
				Max: 11,
			},
			to: Range{
				Min: 1,
				Max: 11,
			},
			want: 6,
		},
		{
			name: "different ranges",
			x:    0.5,
			from: Range{
				Min: 0,
				Max: 1,
			},
			to: Range{
				Min: -1,
				Max: 1,
			},
			want: 0,
		},
		{
			name: "negative limits work",
			x:    2,
			from: Range{
				Min: 0,
				Max: 10,
			},
			to: Range{
				Min: -2,
				Max: 18,
			},
			want: 2,
		},
		{
			name: "negative ranges work",
			x:    5,
			from: Range{
				Min: 0,
				Max: 10,
			},
			to: Range{
				Min: -20,
				Max: 0,
			},
			want: -10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Transpose(tt.x, tt.from, tt.to); got != tt.want {
				t.Errorf("Transpose() = %v, want %v", got, tt.want)
			}
		})
	}
}
