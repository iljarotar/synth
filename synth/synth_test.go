package synth

import "testing"

func Test_secondsToStep(t *testing.T) {
	tests := []struct {
		name       string
		seconds    float64
		delta      float64
		sampleRate float64
		want       float64
	}{
		{
			name:       "when seconds is 0 step is delta",
			seconds:    0,
			delta:      1,
			sampleRate: 1,
			want:       1,
		},
		{
			name:       "when seconds is 1 step is delta/sampleRate",
			seconds:    1,
			delta:      0.5,
			sampleRate: 1000,
			want:       0.5 / 1000,
		},
		{
			name:       "when seconds is greater than 1 step is delta/(seconds*sampleRate)",
			seconds:    5,
			delta:      0.1,
			sampleRate: 1000,
			want:       0.1 / 5000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := secondsToStep(tt.seconds, tt.delta, tt.sampleRate); got != tt.want {
				t.Errorf("secondsToStep() = %v, want %v", got, tt.want)
			}
		})
	}
}
