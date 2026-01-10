package module

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSampler_Step(t *testing.T) {
	tests := []struct {
		name        string
		s           *Sampler
		modules     ModuleMap
		want        float64
		wantTrigger float64
	}{
		{
			name: "trigger zero",
			s: &Sampler{
				Module: Module{
					current: Output{
						Mono: 0.5,
					},
				},
				In:           "in",
				Trigger:      "trigger",
				triggerValue: 1,
			},
			modules: ModuleMap{
				"in": &Module{
					current: Output{
						Mono: 1,
					},
				},
				"trigger": &Module{
					current: Output{
						Mono: 0,
					},
				},
			},
			want:        0.5,
			wantTrigger: 0,
		},
		{
			name: "trigger transitions to positive",
			s: &Sampler{
				Module: Module{
					current: Output{
						Mono: 0,
					},
				},
				In:           "in",
				Trigger:      "trigger",
				triggerValue: -1,
			},
			modules: ModuleMap{
				"in": &Module{
					current: Output{
						Mono: 1,
					},
				},
				"trigger": &Module{
					current: Output{
						Mono: 1,
					},
				},
			},
			want:        1,
			wantTrigger: 1,
		},
		{
			name: "trigger stays positive",
			s: &Sampler{
				Module: Module{
					current: Output{
						Mono: 0.5,
					},
				},
				In:           "in",
				Trigger:      "trigger",
				triggerValue: 1,
			},
			modules: ModuleMap{
				"in": &Module{
					current: Output{
						Mono: 1,
					},
				},
				"trigger": &Module{
					current: Output{
						Mono: 0.5,
					},
				},
			},
			want:        0.5,
			wantTrigger: 0.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Step(tt.modules)

			if tt.s.current.Mono != tt.want {
				t.Errorf("Sampler.Step() = %v, want %v", tt.s.current.Mono, tt.want)
			}
			if tt.s.triggerValue != tt.wantTrigger {
				t.Errorf("Sampler.Step() trigger = %v, want %v", tt.s.triggerValue, tt.wantTrigger)
			}
		})
	}
}

func TestSampler_Update(t *testing.T) {
	tests := []struct {
		name string
		s    *Sampler
		new  *Sampler
		want *Sampler
	}{
		{
			name: "no update necessary",
			s: &Sampler{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				In:           "in",
				Trigger:      "trigger",
				triggerValue: 1,
			},
			new: nil,
			want: &Sampler{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				In:           "in",
				Trigger:      "trigger",
				triggerValue: 1,
			},
		},
		{
			name: "update all",
			s: &Sampler{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				In:           "in",
				Trigger:      "trigger",
				triggerValue: 1,
			},
			new: &Sampler{
				In:      "new-in",
				Trigger: "new-trigger",
			},
			want: &Sampler{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				In:           "new-in",
				Trigger:      "new-trigger",
				triggerValue: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Update(tt.new)
			if diff := cmp.Diff(tt.want, tt.s, cmp.AllowUnexported(Module{}, Sampler{})); diff != "" {
				t.Errorf("Sampler.Update() diff = %s", diff)
			}
		})
	}
}
