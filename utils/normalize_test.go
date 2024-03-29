package utils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPercentage(t *testing.T) {
	type args struct {
		x   float64
		min float64
		max float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "minimum is 0",
			args: args{
				x:   100,
				min: 100,
				max: 200,
			},
			want: 0,
		},
		{
			name: "middle is 0.5",
			args: args{
				x:   150,
				min: 100,
				max: 200,
			},
			want: 0.5,
		},
		{
			name: "maximum is 1",
			args: args{
				x:   200,
				min: 100,
				max: 200,
			},
			want: 1,
		},
		{
			name: "higher than max",
			args: args{
				x:   300,
				min: 100,
				max: 200,
			},
			want: 2,
		},
		{
			name: "lower than min",
			args: args{
				x:   0,
				min: 100,
				max: 200,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Percentage(tt.args.x, tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("Percentage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLimit(t *testing.T) {
	type args struct {
		x   float64
		min float64
		max float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "idempotent on the lower end",
			args: args{
				x:   100,
				min: 100,
				max: 200,
			},
			want: 100,
		},
		{
			name: "idempotent on the upper end",
			args: args{
				x:   200,
				min: 100,
				max: 200,
			},
			want: 200,
		},
		{
			name: "upper limit",
			args: args{
				x:   300,
				min: 100,
				max: 200,
			},
			want: 200,
		},
		{
			name: "lower limit",
			args: args{
				x:   0,
				min: 100,
				max: 200,
			},
			want: 100,
		},
		{
			name: "don't limit",
			args: args{
				x:   0.1,
				min: 0,
				max: 1,
			},
			want: 0.1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Limit(tt.args.x, tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	type args struct {
		signal []float64
		min    float64
		max    float64
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			name: "unchanged",
			args: args{
				signal: []float64{1, 0, -1},
				min:    -1,
				max:    1,
			},
			want: []float64{1, 0, -1},
		},
		{
			name: "symmetric",
			args: args{
				signal: []float64{2, 0, -2},
				min:    -1,
				max:    1,
			},
			want: []float64{1, 0, -1},
		},
		{
			name: "asymmetric",
			args: args{
				signal: []float64{1, 0, -3},
				min:    -1,
				max:    1,
			},
			want: []float64{1, 0.5, -1},
		},
		{
			name: "asymmetric range",
			args: args{
				signal: []float64{1, 0, -3},
				min:    0,
				max:    2,
			},
			want: []float64{2, 1.5, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Normalize(tt.args.signal, tt.args.min, tt.args.max); !cmp.Equal(got, tt.want) {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}