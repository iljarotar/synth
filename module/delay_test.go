package module

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDelay_Update(t *testing.T) {
	tests := []struct {
		name string
		d    *Delay
		new  *Delay
		want *Delay
	}{
		{
			name: "no update necessary",
			d: &Delay{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Time:       30,
				Mix:        0.5,
				In:         "in",
				CV:         "cv",
				Mod:        "mod",
				Fade:       1,
				sampleRate: 44100,
				c: &comb{
					y:          []float64{0.5, 0, 0, 0.25},
					sampleRate: 44100,
				},
				mixFader: &fader{
					current: 0.5,
					target:  0.5,
					step:    0.1,
				},
			},
			new: nil,
			want: &Delay{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Time:       30,
				Mix:        0.5,
				In:         "in",
				CV:         "cv",
				Mod:        "mod",
				Fade:       1,
				sampleRate: 44100,
				c: &comb{
					y:          []float64{0.5, 0, 0, 0.25},
					sampleRate: 44100,
				},
				mixFader: &fader{
					current: 0.5,
					target:  0.5,
					step:    0.1,
				},
			},
		},
		{
			name: "update all",
			d: &Delay{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Time:       900,
				Mix:        0.5,
				In:         "in",
				CV:         "cv",
				Mod:        "mod",
				Fade:       1,
				sampleRate: 6,
				c: &comb{
					y:          []float64{0.5, 0, 0, 0, 0.25},
					sampleRate: 6,
				},
				mixFader: &fader{
					current: 0.5,
					target:  0.5,
					step:    0.1,
				},
			},
			new: &Delay{
				Time: 1000,
				Mix:  0.25,
				In:   "new-in",
				CV:   "new-cv",
				Mod:  "new-mod",
				Fade: 2,
			},
			want: &Delay{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Time:       1000,
				Mix:        0.5,
				In:         "new-in",
				CV:         "new-cv",
				Mod:        "new-mod",
				Fade:       2,
				sampleRate: 6,
				c: &comb{
					y:          []float64{0.5, 0, 0, 0, 0.25, 0},
					sampleRate: 6,
				},
				mixFader: &fader{
					current: 0.5,
					target:  0.25,
					step:    -0.125 / 6,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.Update(tt.new)
			if diff := cmp.Diff(tt.want, tt.d, cmp.AllowUnexported(Module{}, Delay{}, fader{}, comb{})); diff != "" {
				t.Errorf("Delay.Update() diff = %s", diff)
			}
		})
	}
}

func TestDelay_Step(t *testing.T) {
	tests := []struct {
		name    string
		modules *ModuleMap
		d       *Delay
		want    float64
	}{
		{
			name: "time zero",
			modules: NewModuleMap(map[string]IModule{
				"in": &Module{
					current: Output{
						Mono: 1,
					},
				},
			}),
			d: &Delay{
				Mix: 0.25,
				In:  "in",
				c: &comb{
					y:   []float64{},
					idx: 2,
				},
			},
			want: 0.75,
		},
		{
			name: "non-zero time",
			modules: NewModuleMap(map[string]IModule{
				"in": &Module{
					current: Output{
						Mono: 1,
					},
				},
			}),
			d: &Delay{
				Mix: 0.25,
				In:  "in",
				c: &comb{
					y:   []float64{0.5, 0.25, 0, 0, 0, 0},
					idx: 1,
				},
			},
			want: 0.8125,
		},
		{
			name: "cv",
			modules: NewModuleMap(map[string]IModule{
				"in": &Module{
					current: Output{
						Mono: 1,
					},
				},
				"cv": &Module{
					current: Output{
						Mono: 0.5,
					},
				},
			}),
			d: &Delay{
				Mix: 0.25,
				In:  "in",
				CV:  "cv",
				c: &comb{
					y:   []float64{0.5, 0.25, 0, 0, 0, 0},
					idx: 1,
				},
			},
			want: 0.625,
		},
		{
			name: "mod and cv",
			modules: NewModuleMap(map[string]IModule{
				"in": &Module{
					current: Output{
						Mono: 1,
					},
				},
				"cv": &Module{
					current: Output{
						Mono: 0.5,
					},
				},
				"mod": &Module{
					current: Output{
						Mono: -0.2,
					},
				},
			}),
			d: &Delay{
				Mix: 0.25,
				In:  "in",
				CV:  "cv",
				Mod: "mod",
				c: &comb{
					y:   []float64{0.5, 0.25, 0, 0, 0, 0},
					idx: 1,
				},
			},
			want: 0.7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.Step(tt.modules)
			if tt.d.current.Mono != tt.want {
				t.Errorf("Delay.Step = %v, want %v", tt.d.current.Mono, tt.want)
			}
		})
	}
}

func TestDelay_fade(t *testing.T) {
	tests := []struct {
		name string
		d    *Delay
		want *Delay
	}{
		{
			name: "no fade necessary",
			d: &Delay{
				Mix: 0.5,
				mixFader: &fader{
					current: 0.5,
					target:  0.5,
					step:    0.2,
				},
			},
			want: &Delay{
				Mix: 0.5,
				mixFader: &fader{
					current: 0.5,
					target:  0.5,
				},
			},
		},
		{
			name: "fade",
			d: &Delay{
				Mix: 0.5,
				mixFader: &fader{
					current: 0.5,
					target:  0.4,
					step:    -0.05,
				},
			},
			want: &Delay{
				Mix: 0.45,
				mixFader: &fader{
					current: 0.45,
					target:  0.4,
					step:    -0.05,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.fade()
			if diff := cmp.Diff(tt.want, tt.d, cmp.AllowUnexported(Module{}, Delay{}, fader{})); diff != "" {
				t.Errorf("Delay.fade() diff = %s", diff)
			}
		})
	}
}
