package module

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/iljarotar/synth/concurrency"
)

func TestMixer_Step(t *testing.T) {
	tests := []struct {
		name    string
		m       *Mixer
		modules *ModuleMap
		want    Output
	}{
		{
			name: "no input",
			m: &Mixer{
				Gain:        0,
				In:          map[string]float64{},
				sampleRate:  1,
				inputFaders: concurrency.NewSyncMap(map[string]*fader{}),
				in:          concurrency.NewSyncMap(map[string]float64{}),
			},
			modules: NewModuleMap(map[string]IModule{
				"in": &Module{},
			}),
			want: Output{},
		},
		{
			name: "no module found",
			m: &Mixer{
				Gain: 1,
				In: map[string]float64{
					"sine": 1,
				},
				sampleRate:  1,
				inputFaders: concurrency.NewSyncMap(map[string]*fader{}),
				in:          concurrency.NewSyncMap(map[string]float64{}),
			},
			modules: NewModuleMap(map[string]IModule{
				"in": &Oscillator{
					Module: Module{},
				},
			}),
			want: Output{},
		},
		{
			name: "input gain 0",
			m: &Mixer{
				Gain: 1,
				In: map[string]float64{
					"in": 0,
				},
				sampleRate:  1,
				inputFaders: concurrency.NewSyncMap(map[string]*fader{}),
				in:          concurrency.NewSyncMap(map[string]float64{}),
			},
			modules: NewModuleMap(map[string]IModule{
				"in": &Module{
					current: Output{
						Mono:  1,
						Left:  0.5,
						Right: 0.5,
					},
				},
			}),
			want: Output{},
		},
		{
			name: "input",
			m: &Mixer{
				Gain: 1,
				In: map[string]float64{
					"in": 1,
				},
				sampleRate:  1,
				inputFaders: concurrency.NewSyncMap(map[string]*fader{}),
				in: concurrency.NewSyncMap(map[string]float64{
					"in": 1,
				}),
			},
			modules: NewModuleMap(map[string]IModule{
				"in": &Module{
					current: Output{
						Mono:  1,
						Left:  0.5,
						Right: 0.5,
					},
				},
			}),
			want: Output{
				Mono:  1,
				Left:  0.5,
				Right: 0.5,
			},
		},
		{
			name: "mod",
			m: &Mixer{
				Gain: 0.5,
				Mod:  "lfo",
				In: map[string]float64{
					"in": 1,
				},
				sampleRate:  1,
				inputFaders: concurrency.NewSyncMap(map[string]*fader{}),
				in: concurrency.NewSyncMap(map[string]float64{
					"in": 1,
				}),
			},
			modules: NewModuleMap(map[string]IModule{
				"in": &Module{
					current: Output{
						Mono:  1,
						Left:  0.5,
						Right: 0.5,
					},
				},
				"lfo": &Module{
					current: Output{
						Mono:  0.5,
						Left:  0.25,
						Right: 0.25,
					},
				},
			}),
			want: Output{
				Mono:  1 * 0.75,
				Left:  0.5 * 0.75,
				Right: 0.5 * 0.75,
			},
		},
		{
			name: "cv",
			m: &Mixer{
				Gain: 0.5,
				CV:   "cv",
				In: map[string]float64{
					"in": 1,
				},
				sampleRate:  1,
				inputFaders: concurrency.NewSyncMap(map[string]*fader{}),
				in: concurrency.NewSyncMap(map[string]float64{
					"in": 1,
				}),
			},
			modules: NewModuleMap(map[string]IModule{
				"in": &Module{
					current: Output{
						Mono:  1,
						Left:  0.5,
						Right: 0.5,
					},
				},
				"cv": &Module{
					current: Output{
						Mono:  0.5,
						Left:  0.25,
						Right: 0.25,
					},
				},
			}),
			want: Output{
				Mono:  1 * 0.75,
				Left:  0.5 * 0.75,
				Right: 0.5 * 0.75,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Step(tt.modules)

			if diff := cmp.Diff(tt.want, tt.m.Current()); diff != "" {
				t.Errorf("Mixer.Step() diff = %s", diff)
			}
		})
	}
}

func TestMixer_Update(t *testing.T) {
	sampleRate := 44100.0

	tests := []struct {
		name string
		m    *Mixer
		new  *Mixer
		want *Mixer
	}{
		{
			name: "no update necessary",
			m: &Mixer{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Gain: 1,
				CV:   "cv",
				Mod:  "mod",
				In: map[string]float64{
					"in1": 1,
				},
				Fade:       1,
				sampleRate: sampleRate,
				gainFader: &fader{
					current: 1,
					target:  1,
					step:    0.5,
				},
				inputFaders: concurrency.NewSyncMap(map[string]*fader{
					"in1": {
						current: 1,
						target:  1,
						step:    0.5,
					},
				}),
				in: concurrency.NewSyncMap(map[string]float64{
					"in": 1,
				}),
			},
			new: nil,
			want: &Mixer{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Gain: 1,
				CV:   "cv",
				Mod:  "mod",
				In: map[string]float64{
					"in1": 1,
				},
				Fade:       1,
				sampleRate: sampleRate,
				gainFader: &fader{
					current: 1,
					target:  1,
					step:    0.5,
				},
				inputFaders: concurrency.NewSyncMap(map[string]*fader{
					"in1": {
						current: 1,
						target:  1,
						step:    0.5,
					},
				}),
				in: concurrency.NewSyncMap(map[string]float64{
					"in": 1,
				}),
			},
		},
		{
			name: "update all",
			m: &Mixer{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Gain: 1,
				CV:   "cv",
				Mod:  "mod",
				In: map[string]float64{
					"in1": 1,
					"in2": 1,
				},
				Fade:       1,
				sampleRate: sampleRate,
				gainFader: &fader{
					current: 1,
					target:  1,
					step:    0.5,
				},
				inputFaders: concurrency.NewSyncMap(map[string]*fader{
					"in1": {
						current: 1,
						target:  1,
						step:    0.5,
					},
					"in2": {
						current: 1,
						target:  1,
						step:    0.5,
					},
				}),
				in: concurrency.NewSyncMap(map[string]float64{
					"in1": 1,
					"in2": 1,
				}),
			},
			new: &Mixer{
				Gain: 0.5,
				CV:   "new-cv",
				Mod:  "new-mod",
				In: map[string]float64{
					"in1": 0.5,
					"in3": 0.5,
				},
				Fade: 2,
				in: concurrency.NewSyncMap(map[string]float64{
					"in1": 0.5,
					"in3": 0.5,
				}),
			},
			want: &Mixer{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Gain: 1,
				CV:   "new-cv",
				Mod:  "new-mod",
				In: map[string]float64{
					"in1": 0.5,
					"in3": 0.5,
				},
				Fade:       2,
				sampleRate: sampleRate,
				gainFader: &fader{
					current: 1,
					target:  0.5,
					step:    -0.25 / sampleRate,
				},
				inputFaders: concurrency.NewSyncMap(map[string]*fader{
					"in1": {
						current: 1,
						target:  0.5,
						step:    -0.25 / sampleRate,
					},
					"in2": {
						current: 1,
						target:  0,
						step:    -0.5 / sampleRate,
					},
					"in3": {
						current: 0,
						target:  0.5,
						step:    0.25 / sampleRate,
					},
				}),
				in: concurrency.NewSyncMap(map[string]float64{
					"in1": 1,
					"in2": 1,
					"in3": 0,
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Update(tt.new)
			if diff := cmp.Diff(
				tt.want,
				tt.m,
				cmp.AllowUnexported(
					Module{},
					Mixer{},
					fader{},
				),
				cmpopts.IgnoreUnexported(concurrency.SyncMap[string, *fader]{}, concurrency.SyncMap[string, float64]{}),
			); diff != "" {
				t.Errorf("Mixer.Update() diff = %s", diff)
			}
		})
	}
}

func TestMixer_fade(t *testing.T) {
	tests := []struct {
		name string
		m    *Mixer
		want *Mixer
	}{
		{
			name: "no fade necessary",
			m: &Mixer{
				Gain: 1,
				In: map[string]float64{
					"in1": 1,
					"in2": 1,
				},
				gainFader: &fader{
					current: 1,
					target:  1,
					step:    0.5,
				},
				inputFaders: concurrency.NewSyncMap(map[string]*fader{
					"in1": {
						current: 1,
						target:  1,
						step:    0.5,
					},
					"in2": {
						current: 1,
						target:  1,
						step:    0.5,
					},
				}),
				in: concurrency.NewSyncMap(map[string]float64{
					"in1": 1,
					"in2": 1,
				}),
			},
			want: &Mixer{
				Gain: 1,
				In: map[string]float64{
					"in1": 1,
					"in2": 1,
				},
				gainFader: &fader{
					current: 1,
					target:  1,
				},
				inputFaders: concurrency.NewSyncMap(map[string]*fader{
					"in1": {
						current: 1,
						target:  1,
					},
					"in2": {
						current: 1,
						target:  1,
					},
				}),
				in: concurrency.NewSyncMap(map[string]float64{
					"in1": 1,
					"in2": 1,
				}),
			},
		},
		{
			name: "fade all",
			m: &Mixer{
				Gain: 1,
				In: map[string]float64{
					"in1": 1,
					"in2": 0.1,
				},
				gainFader: &fader{
					current: 1,
					target:  0.5,
					step:    -0.1,
				},
				inputFaders: concurrency.NewSyncMap(map[string]*fader{
					"in1": {
						current: 1,
						target:  0.5,
						step:    -0.2,
					},
					"in2": {
						current: 0.1,
						target:  0,
						step:    -0.2,
					},
				}),
				in: concurrency.NewSyncMap(map[string]float64{
					"in1": 1,
					"in2": 0.1,
				}),
			},
			want: &Mixer{
				Gain: 0.9,
				In: map[string]float64{
					"in1": 1,
					"in2": 0.1,
				},
				gainFader: &fader{
					current: 0.9,
					target:  0.5,
					step:    -0.1,
				},
				inputFaders: concurrency.NewSyncMap(map[string]*fader{
					"in1": {
						current: 0.8,
						target:  0.5,
						step:    -0.2,
					},
				}),
				in: concurrency.NewSyncMap(map[string]float64{
					"in1": 0.8,
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.fade()
			if diff := cmp.Diff(
				tt.want,
				tt.m,
				cmp.AllowUnexported(
					Module{},
					Mixer{},
					fader{},
				),
				cmpopts.IgnoreUnexported(concurrency.SyncMap[string, *fader]{}, concurrency.SyncMap[string, float64]{}),
			); diff != "" {
				t.Errorf("Mixer.fade() diff = %s", diff)
			}
		})
	}
}
