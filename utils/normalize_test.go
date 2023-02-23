package utils

import (
	"testing"
)

func TestInverseNormalize(t *testing.T) {
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
			name: "0 lies in the middle",
			args: args{
				x:   0,
				min: 100,
				max: 200,
			},
			want: 150,
		},
		{
			name: "1 is maximum",
			args: args{
				x:   1,
				min: 100,
				max: 200,
			},
			want: 200,
		},
		{
			name: "-1 is minimum",
			args: args{
				x:   -1,
				min: 100,
				max: 200,
			},
			want: 100,
		},
		{
			name: "exceeding 1",
			args: args{
				x:   3,
				min: 100,
				max: 200,
			},
			want: 300,
		},
		{
			name: "lower than -1",
			args: args{
				x:   -3,
				min: 100,
				max: 200,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InverseNormalize(tt.args.x, tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			name: "lower is -1",
			args: args{
				x:   100,
				min: 100,
				max: 200,
			},
			want: -1,
		},
		{
			name: "upper is 1",
			args: args{
				x:   200,
				min: 100,
				max: 200,
			},
			want: 1,
		},
		{
			name: "middle is 0",
			args: args{
				x:   150,
				min: 100,
				max: 200,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Normalize(tt.args.x, tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranspose(t *testing.T) {
	type args struct {
		x      float64
		oldMin float64
		oldMax float64
		newMin float64
		newMax float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "lower bound",
			args: args{
				x:      10,
				oldMin: 10,
				oldMax: 20,
				newMin: -2,
				newMax: 2,
			},
			want: -2,
		},
		{
			name: "upper bound",
			args: args{
				x:      20,
				oldMin: 10,
				oldMax: 20,
				newMin: -2,
				newMax: 2,
			},
			want: 2,
		},
		{
			name: "middle",
			args: args{
				x:      15,
				oldMin: 10,
				oldMax: 20,
				newMin: -2,
				newMax: 2,
			},
			want: 0,
		},
		{
			name: "negative to 0",
			args: args{
				x:      -6,
				oldMin: -6,
				oldMax: 0,
				newMin: -2,
				newMax: 2,
			},
			want: -2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Transpose(tt.args.x, tt.args.oldMin, tt.args.oldMax, tt.args.newMin, tt.args.newMax); got != tt.want {
				t.Errorf("Transpose() = %v, want %v", got, tt.want)
			}
		})
	}
}
