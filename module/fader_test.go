package module

import (
	"testing"
)

func Test_fader_initialize(t *testing.T) {
	tests := []struct {
		name       string
		duration   float64
		sampleRate float64
		f          *fader
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.initialize(tt.duration, tt.sampleRate)
		})
	}
}

func Test_fader_fade(t *testing.T) {
	tests := []struct {
		name string
		f    *fader
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.fade()
		})
	}
}
