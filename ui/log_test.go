package ui

import "testing"

func Test_formatTime(t *testing.T) {
	tests := []struct {
		name string
		time int
		want string
	}{
		{
			name: "less than 10 seconds",
			time: 9,
			want: "00:00:09",
		},
		{
			name: "between 10 seconds and 1 minute",
			time: 59,
			want: "00:00:59",
		},
		{
			name: "between 1 minute and 10 minutes",
			time: 599,
			want: "00:09:59",
		},
		{
			name: "between 10 minutes and 1 hour",
			time: 3599,
			want: "00:59:59",
		},
		{
			name: "more than 1 hour",
			time: 3731,
			want: "01:02:11",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatTime(tt.time); got != tt.want {
				t.Errorf("formatTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
