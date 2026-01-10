package module

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFilter_Update(t *testing.T) {
	sampleRate := 44100.0

	tests := []struct {
		name string
		f    *Filter
		new  *Filter
		want *Filter
	}{
		{
			name: "no update necessary",
			f: &Filter{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Type:       "BandPass",
				Freq:       440,
				Width:      50,
				CV:         "cv",
				Mod:        "mod",
				In:         "in",
				Fade:       1,
				sampleRate: sampleRate,
				a0:         1,
				a1:         1,
				a2:         1,
				b0:         1,
				b1:         1,
				b2:         1,
				inputs: filterInputs{
					x0: 1,
					x1: 1,
					x2: 1,
					y0: 1,
					y1: 1,
				},
				freqFader: &fader{
					current: 440,
					target:  440,
					step:    0,
				},
				widthFader: &fader{
					current: 50,
					target:  50,
					step:    0,
				},
			},
			new: nil,
			want: &Filter{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Type:       "BandPass",
				Freq:       440,
				Width:      50,
				CV:         "cv",
				Mod:        "mod",
				In:         "in",
				Fade:       1,
				sampleRate: sampleRate,
				a0:         1,
				a1:         1,
				a2:         1,
				b0:         1,
				b1:         1,
				b2:         1,
				inputs: filterInputs{
					x0: 1,
					x1: 1,
					x2: 1,
					y0: 1,
					y1: 1,
				},
				freqFader: &fader{
					current: 440,
					target:  440,
					step:    0,
				},
				widthFader: &fader{
					current: 50,
					target:  50,
					step:    0,
				},
			},
		},
		{
			name: "update all",
			f: &Filter{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Type:       "BandPass",
				Freq:       440,
				Width:      50,
				CV:         "cv",
				Mod:        "mod",
				In:         "in",
				Fade:       1,
				sampleRate: sampleRate,
				a0:         1,
				a1:         1,
				a2:         1,
				b0:         1,
				b1:         1,
				b2:         1,
				inputs: filterInputs{
					x0: 1,
					x1: 1,
					x2: 1,
					y0: 1,
					y1: 1,
				},
				freqFader: &fader{
					current: 440,
					target:  440,
					step:    0,
				},
				widthFader: &fader{
					current: 50,
					target:  50,
					step:    0,
				},
			},
			new: &Filter{
				Type:  "LowPass",
				Freq:  220,
				Width: 0,
				CV:    "new-cv",
				Mod:   "new-mod",
				In:    "new-in",
				Fade:  2,
				a0:    2,
				a1:    2,
				a2:    2,
				b0:    2,
				b1:    2,
				b2:    2,
			},
			want: &Filter{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Type:       "LowPass",
				Freq:       440,
				Width:      50,
				CV:         "new-cv",
				Mod:        "new-mod",
				In:         "new-in",
				Fade:       2,
				sampleRate: sampleRate,
				a0:         2,
				a1:         2,
				a2:         2,
				b0:         2,
				b1:         2,
				b2:         2,
				inputs: filterInputs{
					x0: 1,
					x1: 1,
					x2: 1,
					y0: 1,
					y1: 1,
				},
				freqFader: &fader{
					current: 440,
					target:  220,
					step:    -110 / sampleRate,
				},
				widthFader: &fader{
					current: 50,
					target:  0,
					step:    -25 / sampleRate,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.Update(tt.new)
			if diff := cmp.Diff(tt.want, tt.f, cmp.AllowUnexported(Module{}, Filter{}, fader{}, filterInputs{})); diff != "" {
				t.Errorf("Filter.Update() diff = %s", diff)
			}
		})
	}
}

func TestFilter_fade(t *testing.T) {
	tests := []struct {
		name string
		f    *Filter
		want *Filter
	}{
		{
			name: "no fade necessary",
			f: &Filter{
				Freq:  440,
				Width: 50,
				freqFader: &fader{
					current: 440,
					target:  440,
					step:    25,
				},
				widthFader: &fader{
					current: 50,
					target:  50,
					step:    10,
				},
			},
			want: &Filter{
				Freq:  440,
				Width: 50,
				freqFader: &fader{
					current: 440,
					target:  440,
				},
				widthFader: &fader{
					current: 50,
					target:  50,
				},
			},
		},
		{
			name: "fade all",
			f: &Filter{
				Freq:  440,
				Width: 50,
				freqFader: &fader{
					current: 440,
					target:  220,
					step:    -20,
				},
				widthFader: &fader{
					current: 50,
					target:  0,
					step:    -10,
				},
			},
			want: &Filter{
				Freq:  420,
				Width: 40,
				freqFader: &fader{
					current: 420,
					target:  220,
					step:    -20,
				},
				widthFader: &fader{
					current: 40,
					target:  0,
					step:    -10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.fade()
			if diff := cmp.Diff(tt.want, tt.f, cmp.AllowUnexported(Module{}, Filter{}, fader{}, filterInputs{})); diff != "" {
				t.Errorf("Filter.fade() diff = %s", diff)
			}
		})
	}
}
