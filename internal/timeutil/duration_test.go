package timeutil

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want time.Duration
	}{
		{name: "minutes", in: "30m", want: 30 * time.Minute},
		{name: "hours", in: "3h", want: 3 * time.Hour},
		{name: "days", in: "7d", want: 7 * 24 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.in)
			if err != nil {
				t.Fatalf("ParseDuration() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("ParseDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDurationRejectsInvalidValues(t *testing.T) {
	for _, in := range []string{"0m", "-1h", "10x", "abc"} {
		t.Run(in, func(t *testing.T) {
			if _, err := ParseDuration(in); err == nil {
				t.Fatalf("ParseDuration(%q) expected error", in)
			}
		})
	}
}
