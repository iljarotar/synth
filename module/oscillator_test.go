package module

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOscillator_Step(t *testing.T) {
	sampleRate := 44100.0
	twoPi := 2 * math.Pi

	tests := []struct {
		name    string
		modules *ModuleMap
		o       *Oscillator
		want    float64
		wantArg float64
	}{
		{
			name:    "all zero",
			modules: &ModuleMap{},
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
			modules: &ModuleMap{},
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
			modules: &ModuleMap{},
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
			modules: &ModuleMap{
				modules: map[string]IModule{
					"mod": &Module{
						current: Output{
							Mono: 1,
						},
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
			modules: &ModuleMap{
				modules: map[string]IModule{
					"cv": &Module{
						current: Output{
							Mono: 1,
						},
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
			wantArg: twoPi * freqRange.Max / sampleRate,
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

func TestOscillator_Update(t *testing.T) {
	sampleRate := 44100.0

	tests := []struct {
		name string
		o    *Oscillator
		new  *Oscillator
		want *Oscillator
	}{
		{
			name: "no update necessary",
			o: &Oscillator{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Type:       "Sine",
				Freq:       440,
				CV:         "cv",
				Mod:        "mod",
				Phase:      0.5,
				Fade:       1,
				sampleRate: sampleRate,
				arg:        2,
				freqFader: &fader{
					current: 440,
					target:  440,
					step:    12,
				},
				phaseFader: &fader{
					current: 0.5,
					target:  0.5,
					step:    0.1,
				},
			},
			new: nil,
			want: &Oscillator{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Type:       "Sine",
				Freq:       440,
				CV:         "cv",
				Mod:        "mod",
				Phase:      0.5,
				Fade:       1,
				sampleRate: sampleRate,
				arg:        2,
				freqFader: &fader{
					current: 440,
					target:  440,
					step:    12,
				},
				phaseFader: &fader{
					current: 0.5,
					target:  0.5,
					step:    0.1,
				},
			},
		},
		{
			name: "update all",
			o: &Oscillator{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Type:       "Sine",
				Freq:       440,
				CV:         "cv",
				Mod:        "mod",
				Phase:      0.5,
				Fade:       1,
				sampleRate: sampleRate,
				arg:        2,
				freqFader: &fader{
					current: 440,
					target:  440,
					step:    12,
				},
				phaseFader: &fader{
					current: 0.5,
					target:  0.5,
					step:    0.1,
				},
			},
			new: &Oscillator{
				Type:  "Square",
				Freq:  220,
				CV:    "new-cv",
				Mod:   "new-mod",
				Phase: 0,
				Fade:  2,
			},
			want: &Oscillator{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Type:       "Square",
				Freq:       440,
				CV:         "new-cv",
				Mod:        "new-mod",
				Phase:      0.5,
				Fade:       2,
				sampleRate: sampleRate,
				arg:        2,
				freqFader: &fader{
					current: 440,
					target:  220,
					step:    -110 / sampleRate,
				},
				phaseFader: &fader{
					current: 0.5,
					target:  0,
					step:    -0.25 / sampleRate,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.o.Update(tt.new)
			if diff := cmp.Diff(tt.want, tt.o, cmp.AllowUnexported(Module{}, Oscillator{}, fader{})); diff != "" {
				t.Errorf("Oscillator.Update() diff = %s", diff)
			}
		})
	}
}

func TestOscillator_fade(t *testing.T) {
	tests := []struct {
		name string
		o    *Oscillator
		want *Oscillator
	}{
		{
			name: "no fade necessary",
			o: &Oscillator{
				Freq:  440,
				Phase: 0.5,
				freqFader: &fader{
					current: 440,
					target:  440,
					step:    22,
				},
				phaseFader: &fader{
					current: 0.5,
					target:  0.5,
					step:    0.5,
				},
			},
			want: &Oscillator{
				Freq:  440,
				Phase: 0.5,
				freqFader: &fader{
					current: 440,
					target:  440,
				},
				phaseFader: &fader{
					current: 0.5,
					target:  0.5,
				},
			},
		},
		{
			name: "fade all",
			o: &Oscillator{
				Freq:  440,
				Phase: 0.5,
				freqFader: &fader{
					current: 440,
					target:  800,
					step:    20,
				},
				phaseFader: &fader{
					current: 0.5,
					target:  0.75,
					step:    0.1,
				},
			},
			want: &Oscillator{
				Freq:  460,
				Phase: 0.6,
				freqFader: &fader{
					current: 460,
					target:  800,
					step:    20,
				},
				phaseFader: &fader{
					current: 0.6,
					target:  0.75,
					step:    0.1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.o.fade()
			if diff := cmp.Diff(tt.want, tt.o, cmp.AllowUnexported(Module{}, Oscillator{}, fader{})); diff != "" {
				t.Errorf("Oscillator.fade() diff = %s", diff)
			}
		})
	}
}
