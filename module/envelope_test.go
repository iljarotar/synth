package module

import "testing"

func TestEnvelope_trigger(t *testing.T) {
	tests := []struct {
		name     string
		envelope Envelope
		t        float64
		want     *float64
	}{
		{
			name: "no delay must trigger immediately",
			envelope: Envelope{
				currentBPM: 20,
			},
			t:    0,
			want: pointer(0.0),
		},
		{
			name: "before time reaches delay no trigger occurs",
			envelope: Envelope{
				currentBPM: 20,
				Delay:      10,
			},
			t:    5,
			want: nil,
		},
		{
			name: "in between two beats no trigger occurs",
			envelope: Envelope{
				currentBPM:      20,
				Delay:           0,
				lastTriggeredAt: pointer(2.9),
			},
			t:    5,
			want: pointer(2.9),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.envelope.trigger(tt.t, make(ModulesMap))
			got := tt.envelope.lastTriggeredAt
			if got == nil && tt.want == nil {
				return
			}
			if got == nil && tt.want != nil {
				t.Errorf("Envelope.trigger() got %v, want %v", got, *tt.want)
				return
			}
			if got != nil && tt.want == nil {
				t.Errorf("Envelope.trigger() got %v, want %v", *got, tt.want)
				return
			}
			if *got != *tt.want {
				t.Errorf("Envelope.trigger() got %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestEnvelope_getCurrentValue(t *testing.T) {
	tests := []struct {
		name     string
		envelope Envelope
		t        float64
		want     float64
	}{
		{
			name: "value at the beginning of envelope is 0",
			envelope: Envelope{
				currentConfig: envelopeConfig{
					attack:       1,
					decay:        1,
					sustain:      1,
					release:      1,
					peak:         1,
					sustainLevel: 0.5,
				},
				lastTriggeredAt: pointer(0.0),
			},
			t:    0,
			want: 0,
		},
		{
			name: "value at the end of envelope is 0",
			envelope: Envelope{
				currentConfig: envelopeConfig{
					attack:       1,
					decay:        1,
					sustain:      1,
					release:      1,
					peak:         1,
					sustainLevel: 0.5,
				},
				lastTriggeredAt: pointer(0.0),
			},
			t:    4,
			want: 0,
		},
		{
			name: "value at the end of attack is peak",
			envelope: Envelope{
				currentConfig: envelopeConfig{
					attack:       1,
					decay:        1,
					sustain:      1,
					release:      1,
					peak:         1,
					sustainLevel: 0.5,
				},
				lastTriggeredAt: pointer(0.0),
			},
			t:    1,
			want: 1,
		},
		{
			name: "value at the end of decay is sustain level",
			envelope: Envelope{
				currentConfig: envelopeConfig{
					attack:       1,
					decay:        1,
					sustain:      1,
					release:      1,
					peak:         1,
					sustainLevel: 0.5,
				},
				lastTriggeredAt: pointer(0.0),
			},
			t:    2,
			want: 0.5,
		},
		{
			name: "value at the end of sustain is sustain level",
			envelope: Envelope{
				currentConfig: envelopeConfig{
					attack:       1,
					decay:        1,
					sustain:      1,
					release:      1,
					peak:         1,
					sustainLevel: 0.5,
				},
				lastTriggeredAt: pointer(0.0),
			},
			t:    3,
			want: 0.5,
		},
		{
			name: "value in the middle of release is half the sustain level",
			envelope: Envelope{
				currentConfig: envelopeConfig{
					attack:       1,
					decay:        1,
					sustain:      1,
					release:      1,
					peak:         1,
					sustainLevel: 0.5,
				},
				lastTriggeredAt: pointer(0.0),
			},
			t:    3.5,
			want: 0.25,
		},
		{
			name: "while not active value is 0",
			envelope: Envelope{
				currentConfig: envelopeConfig{
					attack:       1,
					decay:        1,
					sustain:      1,
					release:      1,
					peak:         1,
					sustainLevel: 0.5,
				},
				lastTriggeredAt: pointer(4.0),
			},
			t:    12,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.envelope.getCurrentValue(tt.t); got != tt.want {
				t.Errorf("Envelope.getCurrentValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func pointer[T any](value T) *T {
	return &value
}
