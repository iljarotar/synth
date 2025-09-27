package module

import (
	"math"
	"testing"
)

const tolerance = 0.00001

func TestEnvelope_attack(t *testing.T) {
	tests := []struct {
		name string
		t    float64
		e    *Envelope
		want float64
	}{
		{
			name: "just triggered",
			t:    1.5,
			e: &Envelope{
				Attack:      2,
				triggeredAt: 1.5,
			},
			want: 0,
		},
		{
			name: "peak reached",
			t:    3.5,
			e: &Envelope{
				Attack:      2,
				triggeredAt: 1.5,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.attack(tt.t)
			diff := math.Abs(got - tt.want)

			if diff > tolerance {
				t.Errorf("Envelope.attack() = %v, want %v, diff %v", got, tt.want, diff)
			}
		})
	}
}

func TestEnvelope_decay(t *testing.T) {
	tests := []struct {
		name string
		t    float64
		e    *Envelope
		want float64
	}{
		{
			name: "starting decay",
			t:    4,
			e: &Envelope{
				Attack:      1,
				triggeredAt: 3,
			},
			want: 1,
		},
		{
			name: "end of decay",
			t:    5,
			e: &Envelope{
				Attack:      1,
				Decay:       1,
				Level:       0.5,
				triggeredAt: 3,
			},
			want: 0.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.decay(tt.t)
			diff := math.Abs(got - tt.want)

			if diff > tolerance {
				t.Errorf("Envelope.decay() = %v, want %v, diff %v", got, tt.want, diff)
			}
		})
	}
}

func TestEnvelope_release(t *testing.T) {
	tests := []struct {
		name string
		t    float64
		e    *Envelope
		want float64
	}{
		{
			name: "starting release",
			t:    4,
			e: &Envelope{
				Module: Module{
					current: Output{
						Mono: 0.5,
					},
				},
				Release:    2,
				releasedAt: 4,
			},
			want: 0.5,
		},
		{
			name: "end of release",
			t:    6,
			e: &Envelope{
				Module: Module{
					current: Output{
						Mono: 0.5,
					},
				},
				Release:    2,
				releasedAt: 4,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.release(tt.t)
			diff := math.Abs(got - tt.want)

			if diff > tolerance {
				t.Errorf("Envelope.release() = %v, want %v, diff %v", got, tt.want, diff)
			}
		})
	}
}

func TestEnvelope_getValue(t *testing.T) {
	tests := []struct {
		name string
		t    float64
		e    *Envelope
		want float64
	}{
		{
			name: "in between beats",
			t:    5.1,
			e: &Envelope{
				Module: Module{
					current: Output{
						Mono: 0.5,
					},
				},
				Release:     3,
				triggeredAt: 1,
				releasedAt:  2,
			},
			want: 0,
		},
		{
			name: "release",
			t:    2,
			e: &Envelope{
				Module: Module{
					current: Output{
						Mono: 0.5,
					},
				},
				Release:     3,
				triggeredAt: 1,
				releasedAt:  2,
			},
			want: 0.5,
		},
		{
			name: "attack",
			t:    7 - 0.01,
			e: &Envelope{
				Attack:      1,
				triggeredAt: 6,
				releasedAt:  2,
			},
			want: 1,
		},
		{
			name: "decay start",
			t:    7,
			e: &Envelope{
				Attack:      1,
				Decay:       2,
				Level:       0.75,
				triggeredAt: 6,
				releasedAt:  2,
			},
			want: 1,
		},
		{
			name: "decay end",
			t:    9 - 0.01,
			e: &Envelope{
				Attack:      1,
				Decay:       2,
				Level:       0.75,
				triggeredAt: 6,
				releasedAt:  2,
			},
			want: 0.75,
		},
		{
			name: "sustain",
			t:    10,
			e: &Envelope{
				Attack:      1,
				Decay:       2,
				Level:       0.75,
				triggeredAt: 6,
				releasedAt:  2,
			},
			want: 0.75,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.getValue(tt.t)
			diff := math.Abs(got - tt.want)

			if diff > tolerance {
				t.Errorf("Envelope.getValue() = %v, want %v, diff %v", got, tt.want, diff)
			}
		})
	}
}

func TestEnvelope_Step(t *testing.T) {
	tests := []struct {
		name            string
		e               *Envelope
		t               float64
		modules         ModuleMap
		want            float64
		wantTriggeredAt float64
		wantReleasedAt  float64
		wantGateValue   float64
	}{
		{
			name: "first trigger",
			e: &Envelope{
				Module:      Module{},
				Gate:        "gate",
				triggeredAt: 0,
				releasedAt:  0,
				gateValue:   0,
				Attack:      1,
				Decay:       2,
				Release:     2,
				Level:       0.5,
			},
			t: 2,
			modules: ModuleMap{
				"gate": &Module{
					current: Output{
						Mono: 1,
					},
				},
			},
			want:            -1,
			wantTriggeredAt: 2,
			wantGateValue:   1,
		},
		{
			name: "release",
			e: &Envelope{
				Module: Module{
					current: Output{
						Mono: 0.5,
					},
				},
				Gate:        "gate",
				triggeredAt: 2,
				releasedAt:  0,
				gateValue:   1,
				Attack:      1,
				Decay:       2,
				Release:     2,
			},
			t: 5,
			modules: ModuleMap{
				"gate": &Module{
					current: Output{
						Mono: -1,
					},
				},
			},
			want:            0.5,
			wantTriggeredAt: 2,
			wantReleasedAt:  5,
			wantGateValue:   -1,
		},
		{
			name: "noop after release",
			e: &Envelope{
				Module: Module{
					current: Output{
						Mono: 0.5, // can't be, but ensures zero is returned after release has ended.
					},
				},
				Gate:        "gate",
				triggeredAt: 2,
				releasedAt:  5,
				gateValue:   -1,
				Attack:      1,
				Decay:       2,
				Release:     2,
			},
			t: 8,
			modules: ModuleMap{
				"gate": &Module{
					current: Output{
						Mono: -1,
					},
				},
			},
			want:            -1,
			wantTriggeredAt: 2,
			wantReleasedAt:  5,
			wantGateValue:   -1,
		},
		{
			name: "noop during sustain",
			e: &Envelope{
				Module: Module{
					current: Output{
						Mono: 1, // can't be, but ensures level is returned during sustain.
					},
				},
				Gate:        "gate",
				triggeredAt: 2,
				releasedAt:  0,
				gateValue:   1,
				Attack:      1,
				Decay:       2,
				Release:     2,
				Level:       0.75,
			},
			t: 8,
			modules: ModuleMap{
				"gate": &Module{
					current: Output{
						Mono: 1,
					},
				},
			},
			want:            0.5,
			wantTriggeredAt: 2,
			wantReleasedAt:  0,
			wantGateValue:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.Step(tt.t, tt.modules)
			diff := math.Abs(tt.e.current.Mono - tt.want)

			if diff > tolerance {
				t.Errorf("Envelope.Step() = %v, want %v, diff %v", tt.e.current.Mono, tt.want, diff)
			}
			if tt.e.triggeredAt != tt.wantTriggeredAt {
				t.Errorf("Envelope.Step() triggeredAt = %v, want %v", tt.e.triggeredAt, tt.wantTriggeredAt)
			}
			if tt.e.releasedAt != tt.wantReleasedAt {
				t.Errorf("Envelope.Step() releasedAt = %v, want %v", tt.e.releasedAt, tt.wantReleasedAt)
			}
			if tt.e.gateValue != tt.wantGateValue {
				t.Errorf("Envelope.Step() gateValue = %v, want %v", tt.e.gateValue, tt.wantGateValue)
			}
		})
	}
}
