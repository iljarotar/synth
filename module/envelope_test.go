package module

import "testing"

func TestEnvelope_trigger(t *testing.T) {
	tests := []struct {
		name     string
		envelope Envelope
		t        float64
		bpm      float64
		want     float64
	}{
		{
			name:     "no time shift must trigger at 0",
			envelope: Envelope{},
			t:        0,
			bpm:      10,
			want:     0,
		},
		{
			name:     "no time shift must trigger at multiples of seconds between two beats",
			envelope: Envelope{},
			t:        13,
			bpm:      10,
			want:     12,
		},
		{
			name:     "with time shift could first trigger at negative time",
			envelope: Envelope{TimeShift: 10},
			t:        0,
			bpm:      10,
			want:     -2,
		},
		{
			name:     "with time shift must trigger at multiples of seconds between two beats plus time shift",
			envelope: Envelope{TimeShift: 10},
			t:        18,
			bpm:      10,
			want:     16,
		},
		{
			name:     "negative time shift also works",
			envelope: Envelope{TimeShift: -3},
			t:        4,
			bpm:      10,
			want:     3,
		},
		{
			name:     "change from negative triggered time to positive",
			envelope: Envelope{lastTriggeredAt: pointer(-1.0), TimeShift: 5},
			t:        5,
			bpm:      10,
			want:     5,
		},
		{
			name:     "no trigger if t is closer to last trigger than seconds between two beats",
			envelope: Envelope{lastTriggeredAt: pointer(12.0)},
			t:        17,
			bpm:      10,
			want:     12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.envelope.trigger(tt.t, tt.bpm, make(ModulesMap))
			if got := tt.envelope.lastTriggeredAt; *got != tt.want {
				t.Errorf("Envelope.trigger() got %f, want %f", *got, tt.want)
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
