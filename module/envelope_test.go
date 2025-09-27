package module

import (
	"testing"
)

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
			name: "half way",
			t:    2,
			e: &Envelope{
				Attack:      2,
				triggeredAt: 1,
			},
			want: 0.5,
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
			if got != tt.want {
				t.Errorf("Envelope.attack() = %v, want %v", got, tt.want)
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
				Decay:       1,
				triggeredAt: 3,
			},
			want: 1,
		},
		{
			name: "half way",
			t:    4.5,
			e: &Envelope{
				Attack:      1,
				Decay:       1,
				triggeredAt: 3,
			},
			want: 0.5,
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
			if got != tt.want {
				t.Errorf("Envelope.decay() = %v, want %v", got, tt.want)
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
				Release:    2,
				releasedAt: 4,
				level:      0.5,
			},
			want: 0.5,
		},
		{
			name: "half way",
			t:    5,
			e: &Envelope{
				Release:    2,
				releasedAt: 4,
				level:      0.5,
			},
			want: 0.25,
		},
		{
			name: "end of release",
			t:    6,
			e: &Envelope{
				Release:    2,
				releasedAt: 4,
				level:      0.5,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.release(tt.t)
			if got != tt.want {
				t.Errorf("Envelope.release() = %v, want %v", got, tt.want)
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
				Release:     3,
				triggeredAt: 1,
				releasedAt:  2,
				level:       0.5,
			},
			want: 0,
		},
		{
			name: "release half way",
			t:    3.5,
			e: &Envelope{
				Release:     3,
				triggeredAt: 1,
				releasedAt:  2,
				level:       0.5,
			},
			want: 0.25,
		},
		{
			name: "attack half way",
			t:    6.5,
			e: &Envelope{
				Attack:      1,
				triggeredAt: 6,
				releasedAt:  2,
			},
			want: 0.5,
		},
		{
			name: "decay half way",
			t:    8,
			e: &Envelope{
				Attack:      1,
				Decay:       2,
				Level:       0.5,
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
			if got != tt.want {
				t.Errorf("Envelope.getValue() = %v, want %v", got, tt.want)
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
		wantLevel       float64
	}{
		{
			name: "first trigger",
			e: &Envelope{
				Gate:    "gate",
				Attack:  1,
				Decay:   2,
				Release: 2,
				Level:   0.5,
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
				Attack:      1,
				Decay:       2,
				Release:     2,
				triggeredAt: 2,
				gateValue:   1,
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
			wantLevel:       0.75,
		},
		{
			name: "noop after release",
			e: &Envelope{
				Gate:        "gate",
				Attack:      1,
				Decay:       2,
				Release:     2,
				triggeredAt: 2,
				releasedAt:  5,
				gateValue:   -1,
				level:       0.25,
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
			wantLevel:       0.25,
		},
		{
			name: "noop during sustain",
			e: &Envelope{
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
			if tt.e.current.Mono != tt.want {
				t.Errorf("Envelope.Step() = %v, want %v", tt.e.current.Mono, tt.want)
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
			if tt.e.level != tt.wantLevel {
				t.Errorf("Envelope.Step() level = %v, want %v", tt.e.level, tt.wantLevel)
			}
		})
	}
}

func Test_linear(t *testing.T) {
	tests := []struct {
		name        string
		startAt     float64
		endAt       float64
		startValue  float64
		targetValue float64
		t           float64
		want        float64
	}{
		{
			name:        "delta zero",
			startAt:     1,
			endAt:       1,
			startValue:  0,
			targetValue: 1,
			t:           2,
			want:        1,
		},
		{
			name:        "at start",
			startAt:     1,
			endAt:       2,
			startValue:  1,
			targetValue: 0.5,
			t:           1,
			want:        1,
		},
		{
			name:        "middle",
			startAt:     2,
			endAt:       4,
			startValue:  0.5,
			targetValue: 1,
			t:           3,
			want:        0.75,
		},
		{
			name:        "end",
			startAt:     2,
			endAt:       4,
			startValue:  0.5,
			targetValue: 0.2,
			t:           4,
			want:        0.2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := linear(tt.startAt, tt.endAt, tt.startValue, tt.targetValue, tt.t); got != tt.want {
				t.Errorf("linear() = %v, want %v", got, tt.want)
			}
		})
	}
}
