package module

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPan_Step(t *testing.T) {
	tests := []struct {
		name    string
		p       *Pan
		modules *ModuleMap
		want    Output
	}{
		{
			name: "static",
			p: &Pan{
				Pan: 0.5,
				In:  "in",
			},
			modules: NewModuleMap(map[string]IModule{
				"in": &Module{
					current: Output{
						Mono: 1,
					},
				},
			}),
			want: Output{
				Mono:  1,
				Left:  0.25,
				Right: 0.75,
			},
		},
		{
			name: "mod",
			p: &Pan{
				Pan: 0.25,
				Mod: "mod",
				In:  "in",
			},
			modules: NewModuleMap(map[string]IModule{
				"in": &Module{
					current: Output{
						Mono: 1,
					},
				},
				"mod": &Module{
					current: Output{
						Mono: 0.25,
					},
				},
			}),
			want: Output{
				Mono:  1,
				Left:  0.25,
				Right: 0.75,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Step(tt.modules)

			if diff := cmp.Diff(tt.want, tt.p.current); diff != "" {
				t.Errorf("Pan.Step() diff = %s", diff)
			}
		})
	}
}

func TestPan_Update(t *testing.T) {
	sampleRate := 44100.0

	tests := []struct {
		name string
		p    *Pan
		new  *Pan
		want *Pan
	}{
		{
			name: "no udpate necessary",
			p: &Pan{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Pan:        0.5,
				Mod:        "mod",
				In:         "in",
				Fade:       1,
				sampleRate: sampleRate,
				panFader: &fader{
					current: 0.5,
					target:  0.5,
					step:    0.1,
				},
			},
			new: nil,
			want: &Pan{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Pan:        0.5,
				Mod:        "mod",
				In:         "in",
				Fade:       1,
				sampleRate: sampleRate,
				panFader: &fader{
					current: 0.5,
					target:  0.5,
					step:    0.1,
				},
			},
		},
		{
			name: "update all",
			p: &Pan{
				Module: Module{
					current: Output{
						Mono:  1,
						Left:  0.5,
						Right: 0.5,
					},
				},
				Pan:        0.5,
				Mod:        "mod",
				In:         "in",
				Fade:       1,
				sampleRate: sampleRate,
				panFader: &fader{
					current: 0.5,
					target:  0.5,
					step:    0.1,
				},
			},
			new: &Pan{
				Pan:  -0.5,
				Mod:  "new-mod",
				In:   "new-in",
				Fade: 2,
			},
			want: &Pan{
				Module: Module{
					current: Output{
						Mono:  1,
						Left:  0.5,
						Right: 0.5,
					},
				},
				Pan:        0.5,
				Mod:        "new-mod",
				In:         "new-in",
				Fade:       2,
				sampleRate: sampleRate,
				panFader: &fader{
					current: 0.5,
					target:  -0.5,
					step:    -0.5 / sampleRate,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Update(tt.new)
			if diff := cmp.Diff(tt.want, tt.p, cmp.AllowUnexported(Module{}, Pan{}, fader{})); diff != "" {
				t.Errorf("Pan.Update() diff = %s", diff)
			}
		})
	}
}

func TestPan_fade(t *testing.T) {
	tests := []struct {
		name string
		p    *Pan
		want *Pan
	}{
		{
			name: "no fade necessary",
			p: &Pan{
				Pan: 0.5,
				panFader: &fader{
					current: 0.5,
					target:  0.5,
					step:    1,
				},
			},
			want: &Pan{
				Pan: 0.5,
				panFader: &fader{
					current: 0.5,
					target:  0.5,
				},
			},
		},
		{
			name: "fade",
			p: &Pan{
				Pan: 0.5,
				panFader: &fader{
					current: 0.5,
					target:  1,
					step:    0.1,
				},
			},
			want: &Pan{
				Pan: 0.6,
				panFader: &fader{
					current: 0.6,
					target:  1,
					step:    0.1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.fade()
			if diff := cmp.Diff(tt.want, tt.p, cmp.AllowUnexported(Module{}, Pan{}, fader{})); diff != "" {
				t.Errorf("Pan.fade() diff = %s", diff)
			}
		})
	}
}
