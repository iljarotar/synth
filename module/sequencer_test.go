package module

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/iljarotar/synth/calc"
)

func Test_noteToFreq(t *testing.T) {
	tests := []struct {
		name      string
		note      string
		pitch     float64
		transpose float64
		want      float64
		wantErr   bool
	}{
		{
			name:  "pitch",
			note:  "a_4",
			pitch: 441.5,
			want:  441.5,
		},
		{
			name:  "up an octave",
			note:  "a_5",
			pitch: 440,
			want:  880,
		},
		{
			name:  "down an octave",
			note:  "a_3",
			pitch: 440,
			want:  220,
		},
		{
			name:    "invalid suffix",
			note:    "a'",
			pitch:   440,
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid prefix",
			note:    "h_3",
			pitch:   440,
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid octave",
			note:    "a_-1",
			pitch:   440,
			want:    0,
			wantErr: true,
		},
		{
			name:    "octave too high",
			note:    "a_11",
			pitch:   440,
			want:    0,
			wantErr: true,
		},
		{
			name:      "transpose",
			note:      "a_4",
			pitch:     440,
			transpose: 7,
			want:      440 * math.Pow(2, 7/12.),
			wantErr:   false,
		},
		{
			name:  "higher note in same octave",
			note:  "b_4",
			pitch: 440,
			want:  440 * math.Pow(2, 1.0/6),
		},
		{
			name:  "lower note in same octave",
			note:  "c_4",
			pitch: 440,
			want:  440 * math.Pow(2, -9.0/12),
		},
		{
			name:  "higher note in higher octave",
			note:  "a#_6",
			pitch: 440,
			want:  440 * math.Pow(2, 25.0/12),
		},
		{
			name:  "higher note in lower octave",
			note:  "b#_2",
			pitch: 440,
			want:  440 * math.Pow(2, -21.0/12),
		},
		{
			name:  "lower note in lower octave",
			note:  "eb_3",
			pitch: 440,
			want:  440 * math.Pow(2, -18.0/12),
		},
		{
			name:  "lower note in higher octave",
			note:  "f#_5",
			pitch: 440,
			want:  440 * math.Pow(2, 9.0/12),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := noteToFreq(tt.note, tt.pitch, tt.transpose)
			if (err != nil) != tt.wantErr {
				t.Errorf("noteToFreq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("noteToFreq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSequencer_makeSequence(t *testing.T) {
	tests := []struct {
		name    string
		s       *Sequencer
		want    []float64
		wantErr bool
	}{
		{
			name:    "empty sequence",
			s:       &Sequencer{},
			want:    nil,
			wantErr: false,
		},
		{
			name: "sequence with error",
			s: &Sequencer{
				Sequence: []string{"a_4", "bb_3", "e#_11"},
				Pitch:    440,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid sequence",
			s: &Sequencer{
				Sequence: []string{"a_4", "a_3", "a_5"},
				Pitch:    440,
			},
			want:    []float64{440, 220, 880},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.makeSequence(); (err != nil) != tt.wantErr {
				t.Errorf("Sequencer.makeSequence() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, tt.s.sequence); diff != "" {
				t.Errorf("Sequencer.makeSequence() diff = %v", diff)
			}
		})
	}
}

func TestSequencer_Step(t *testing.T) {
	tests := []struct {
		name        string
		s           *Sequencer
		modules     *ModuleMap
		want        float64
		wantTrigger float64
	}{
		{
			name: "before first trigger",
			s: &Sequencer{
				Module:       Module{},
				sequence:     []float64{440, 220, 110},
				idx:          -1,
				triggerValue: 0,
			},
			modules:     &ModuleMap{},
			want:        0,
			wantTrigger: 0,
		},
		{
			name: "next note",
			s: &Sequencer{
				Module:       Module{},
				Trigger:      "trigger",
				sequence:     []float64{440, 220, 110},
				idx:          1,
				triggerValue: 0,
			},
			modules: NewModuleMap(map[string]IModule{
				"trigger": &Module{
					current: Output{
						Mono: 1,
					},
				},
			}),
			want:        calc.Transpose(110, freqRange, cvRange),
			wantTrigger: 1,
		},
		{
			name: "back to start",
			s: &Sequencer{
				Module:       Module{},
				Trigger:      "trigger",
				sequence:     []float64{440, 220, 110},
				idx:          2,
				triggerValue: 0,
			},
			modules: NewModuleMap(map[string]IModule{
				"trigger": &Module{
					current: Output{
						Mono: 1,
					},
				},
			}),
			want:        calc.Transpose(440, freqRange, cvRange),
			wantTrigger: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Step(tt.modules)

			if tt.s.current.Mono != tt.want {
				t.Errorf("Sequencer.Step() = %v, want %v", tt.s.current.Mono, tt.want)
			}
			if tt.s.triggerValue != tt.wantTrigger {
				t.Errorf("Sequencer.Step() triggerValue = %v, want %v", tt.s.triggerValue, tt.wantTrigger)
			}
		})
	}
}

func TestSequencer_Update(t *testing.T) {
	tests := []struct {
		name string
		s    *Sequencer
		new  *Sequencer
		want *Sequencer
	}{
		{
			name: "no update necessary",
			s: &Sequencer{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Sequence:     []string{"a_4"},
				Trigger:      "trigger",
				Pitch:        440,
				Transpose:    1,
				Randomize:    true,
				Index:        1,
				sequence:     []float64{440},
				idx:          1,
				triggerValue: 1,
			},
			new: nil,
			want: &Sequencer{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Sequence:     []string{"a_4"},
				Trigger:      "trigger",
				Pitch:        440,
				Transpose:    1,
				Randomize:    true,
				Index:        1,
				sequence:     []float64{440},
				idx:          1,
				triggerValue: 1,
			},
		},
		{
			name: "update all",
			s: &Sequencer{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Sequence:     []string{"a_4", "a_2", "a_3"},
				Trigger:      "trigger",
				Pitch:        440,
				Transpose:    1,
				Randomize:    true,
				Index:        2,
				sequence:     []float64{440, 110, 220},
				idx:          2,
				triggerValue: 1,
			},
			new: &Sequencer{
				Sequence:  []string{"a_5", "a_3"},
				Trigger:   "new-trigger",
				Pitch:     441,
				Transpose: 2,
				Randomize: false,
				Index:     1,
				sequence:  []float64{880, 220},
			},
			want: &Sequencer{
				Module: Module{
					current: Output{
						Mono: 1,
					},
				},
				Sequence:     []string{"a_5", "a_3"},
				Trigger:      "new-trigger",
				Pitch:        441,
				Transpose:    2,
				Randomize:    false,
				Index:        2,
				sequence:     []float64{880, 220},
				idx:          1,
				triggerValue: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Update(tt.new)
			if diff := cmp.Diff(tt.want, tt.s, cmp.AllowUnexported(Module{}, Sequencer{})); diff != "" {
				t.Errorf("Sequencer.Update() diff = %s", diff)
			}
		})
	}
}

func TestSequencer_initialize(t *testing.T) {
	tests := []struct {
		name    string
		s       *Sequencer
		want    *Sequencer
		wantErr bool
	}{
		{
			name: "set limits correctly",
			s: &Sequencer{
				Sequence:     []string{"a_4", "a_3"},
				Trigger:      "trigger",
				Pitch:        540,
				Transpose:    25,
				Randomize:    true,
				Index:        2,
				sequence:     []float64{},
				idx:          0,
				triggerValue: 0,
			},
			want: &Sequencer{
				Sequence:     []string{"a_4", "a_3"},
				Trigger:      "trigger",
				Pitch:        500,
				Transpose:    24,
				Randomize:    true,
				Index:        1,
				sequence:     []float64{2000, 1000},
				idx:          0,
				triggerValue: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.initialize(); (err != nil) != tt.wantErr {
				t.Errorf("Sequencer.initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, tt.s, cmp.AllowUnexported(Module{}, Sequencer{})); diff != "" {
				t.Errorf("Sequencer.initialize() diff = %s", diff)
			}
		})
	}
}
