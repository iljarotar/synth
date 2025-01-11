package module

import (
	"math"
	"testing"
)

func TestSequence_stringToFreq(t *testing.T) {
	tests := []struct {
		name string
		s    *Sequence
		note string
		want float64
	}{
		{
			name: "pitch",
			s: &Sequence{
				Pitch: 441.5,
			},
			note: "a_4",
			want: 441.5,
		},
		{
			name: "up an octave",
			s: &Sequence{
				Pitch: 440,
			},
			note: "a_5",
			want: 880,
		},
		{
			name: "down an octave",
			s: &Sequence{
				Pitch: 440,
			},
			note: "a_3",
			want: 220,
		},
		{
			name: "invalid suffix",
			s: &Sequence{
				Pitch: 440,
			},
			note: "a'",
			want: 0,
		},
		{
			name: "invalid prefix",
			s: &Sequence{
				Pitch: 440,
			},
			note: "h_3",
			want: 0,
		},
		{
			name: "invalid octave",
			s: &Sequence{
				Pitch: 440,
			},
			note: "a_-1",
			want: 0,
		},
		{
			name: "higher note in same octave",
			s: &Sequence{
				Pitch: 440,
			},
			note: "b_4",
			want: 440 * math.Pow(2, 1.0/6),
		},
		{
			name: "lower note in same octave",
			s: &Sequence{
				Pitch: 440,
			},
			note: "c_4",
			want: 440 * math.Pow(2, -9.0/12),
		},
		{
			name: "higher note in higher octave",
			s: &Sequence{
				Pitch: 440,
			},
			note: "a#_6",
			want: 440 * math.Pow(2, 25.0/12),
		},
		{
			name: "higher note in lower octave",
			s: &Sequence{
				Pitch: 440,
			},
			note: "b#_2",
			want: 440 * math.Pow(2, -21.0/12),
		},
		{
			name: "lower note in lower octave",
			s: &Sequence{
				Pitch: 440,
			},
			note: "eb_3",
			want: 440 * math.Pow(2, -18.0/12),
		},
		{
			name: "lower note in higher octave",
			s: &Sequence{
				Pitch: 440,
			},
			note: "f#_5",
			want: 440 * math.Pow(2, 9.0/12),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.stringToFreq(tt.note); got != tt.want {
				t.Errorf("Sequence.stringToFreq() = %v, want %v", got, tt.want)
			}
		})
	}
}
