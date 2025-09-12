package module

import (
	"math"
	"testing"
)

func TestOscillator_Step(t *testing.T) {
	sampleRate := 44100.0
	twoPi := 2 * math.Pi

	tests := []struct {
		name    string
		modules ModulesMap
		o       *Oscillator
		want    float64
		wantArg float64
	}{
		{
			name:    "all zero",
			modules: ModulesMap{},
			o: &Oscillator{
				Freq:       200,
				signal:     SineSignalFunc(),
				sampleRate: sampleRate,
			},
			want:    0,
			wantArg: twoPi * 200 / sampleRate,
		},
		{
			name:    "after one second",
			modules: ModulesMap{},
			o: &Oscillator{
				Freq:       200,
				signal:     SineSignalFunc(),
				sampleRate: sampleRate,
				Module:     Module{},
				arg:        twoPi * 200,
			},
			want:    math.Sin(twoPi * 200),
			wantArg: twoPi*200 + twoPi*200/sampleRate,
		},
		{
			name:    "phase shift",
			modules: ModulesMap{},
			o: &Oscillator{
				Freq:       200,
				signal:     SineSignalFunc(),
				sampleRate: sampleRate,
				Module:     Module{},
				Phase:      0.75,
			},
			want:    math.Sin(twoPi * 0.75),
			wantArg: twoPi * 200 / sampleRate,
		},
		{
			name: "modulation",
			modules: ModulesMap{
				"mod": &Module{
					current: Output{
						Mono: 1,
					},
				},
			},
			o: &Oscillator{
				Freq:       200,
				signal:     SineSignalFunc(),
				sampleRate: sampleRate,
				Module:     Module{},
				Mod:        "mod",
				arg:        twoPi,
			},
			want:    math.Sin(twoPi),
			wantArg: twoPi + twoPi*200*2/sampleRate,
		},
		{
			name: "cv",
			modules: ModulesMap{
				"cv": &Module{
					current: Output{
						Mono: 1,
					},
				},
			},
			o: &Oscillator{
				Freq:       200,
				signal:     SineSignalFunc(),
				sampleRate: sampleRate,
				Module:     Module{},
				CV:         "cv",
				arg:        0,
			},
			want:    0,
			wantArg: twoPi * freqLimits.Max / sampleRate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.o.Step(tt.modules)

			if tt.o.Current().Mono != tt.want {
				t.Errorf("Oscillator.Step() = %v, want %v", tt.o.Current().Mono, tt.want)
			}
			if tt.o.arg != tt.wantArg {
				t.Errorf("Oscillator.Step() arg = %v, want %v", tt.o.arg, tt.wantArg)
			}
		})
	}
}
