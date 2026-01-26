package module

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWavetable_initialze(t *testing.T) {
	tests := []struct {
		name string
		w    *Wavetable
		want []float64
	}{
		{
			name: "limit exceeding values",
			w: &Wavetable{
				Signal: []float64{-0.75, 0, -2, 0.75, 2},
			},
			want: []float64{-0.75, 0, -1, 0.75, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.initialze(0)

			if diff := cmp.Diff(tt.want, tt.w.Signal); diff != "" {
				t.Errorf("Wavetable.initialize() diff = %s", diff)
			}
		})
	}
}

func TestWavetable_Step(t *testing.T) {
	sampleRate := 44100.0

	tests := []struct {
		name    string
		w       *Wavetable
		modules *ModuleMap
		want    float64
		wantIdx float64
	}{
		{
			name: "no mod no cv",
			w: &Wavetable{
				Freq:       2,
				Signal:     []float64{1, 0, -1, 0},
				sampleRate: sampleRate,
				idx:        44100.0 / 8,
			},
			modules: &ModuleMap{},
			want:    0,
			wantIdx: 44100.0/8 + 8/44100.0,
		},
		{
			name: "cv",
			w: &Wavetable{
				Freq:       2,
				Signal:     []float64{1, 0, -1, 0},
				CV:         "cv",
				sampleRate: sampleRate,
				idx:        0,
			},
			modules: NewModuleMap(map[string]IModule{
				"cv": &Module{
					current: Output{
						Mono: 0,
					},
				},
			}),
			want:    1,
			wantIdx: 4 * freqRange.Max / (2 * sampleRate),
		},
		{
			name: "mod",
			w: &Wavetable{
				Freq:       2,
				Signal:     []float64{1, 0, -1, 0},
				Mod:        "mod",
				sampleRate: sampleRate,
				idx:        2.5,
			},
			modules: NewModuleMap(map[string]IModule{
				"mod": &Module{
					current: Output{
						Mono: 1,
					},
				},
			}),
			want:    -1,
			wantIdx: 2.5 + 16/sampleRate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.Step(tt.modules)
		})
	}
}

func TestWavetable_Update(t *testing.T) {
	sampleRate := 44100.0

	tests := []struct {
		name string
		w    *Wavetable
		new  *Wavetable
		want *Wavetable
	}{
		{
			name: "no update necessary",
			w: &Wavetable{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Freq:       440,
				CV:         "cv",
				Mod:        "mod",
				Signal:     []float64{1, 0, -1, 0},
				Fade:       1,
				sampleRate: sampleRate,
				idx:        1,
				freqFader: &fader{
					current: 440,
					target:  440,
					step:    10,
				},
			},
			new: nil,
			want: &Wavetable{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Freq:       440,
				CV:         "cv",
				Mod:        "mod",
				Signal:     []float64{1, 0, -1, 0},
				Fade:       1,
				sampleRate: sampleRate,
				idx:        1,
				freqFader: &fader{
					current: 440,
					target:  440,
					step:    10,
				},
			},
		},
		{
			name: "update all",
			w: &Wavetable{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Freq:       440,
				CV:         "cv",
				Mod:        "mod",
				Signal:     []float64{1, 0, -1, 0},
				Fade:       1,
				sampleRate: sampleRate,
				idx:        1,
				freqFader: &fader{
					current: 440,
					target:  440,
					step:    10,
				},
			},
			new: &Wavetable{
				Freq:   880,
				CV:     "new-cv",
				Mod:    "new-mod",
				Signal: []float64{0, 1, 0, -1},
				Fade:   2,
			},
			want: &Wavetable{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Freq:       440,
				CV:         "new-cv",
				Mod:        "new-mod",
				Signal:     []float64{0, 1, 0, -1},
				Fade:       2,
				sampleRate: sampleRate,
				idx:        1,
				freqFader: &fader{
					current: 440,
					target:  880,
					step:    220 / sampleRate,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.Update(tt.new)
			if diff := cmp.Diff(tt.want, tt.w, cmp.AllowUnexported(Module{}, Wavetable{}, fader{})); diff != "" {
				t.Errorf("Wavetable.Update() diff = %s", diff)
			}
		})
	}
}

func TestWavetable_fade(t *testing.T) {
	tests := []struct {
		name string
		w    *Wavetable
		want *Wavetable
	}{
		{
			name: "no fade necessary",
			w: &Wavetable{
				Freq: 440,
				freqFader: &fader{
					current: 440,
					target:  440,
					step:    12,
				},
			},
			want: &Wavetable{
				Freq: 440,
				freqFader: &fader{
					current: 440,
					target:  440,
				},
			},
		},
		{
			name: "fade",
			w: &Wavetable{
				Freq: 440,
				freqFader: &fader{
					current: 440,
					target:  800,
					step:    10,
				},
			},
			want: &Wavetable{
				Freq: 450,
				freqFader: &fader{
					current: 450,
					target:  800,
					step:    10},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.fade()
			if diff := cmp.Diff(tt.want, tt.w, cmp.AllowUnexported(Module{}, Wavetable{}, fader{})); diff != "" {
				t.Errorf("Wavetable.fade() diff = %s", diff)
			}
		})
	}
}
