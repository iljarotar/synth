package module

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPan_Step(t *testing.T) {
	tests := []struct {
		name    string
		p       *Pan
		modules ModulesMap
		want    Output
	}{
		{
			name: "static",
			p: &Pan{
				Pan: 0.5,
				In:  "in",
			},
			modules: ModulesMap{
				"in": &Module{
					current: Output{
						Mono: 1,
					},
				},
			},
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
			modules: ModulesMap{
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
			},
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
