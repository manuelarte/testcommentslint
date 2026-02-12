package main

import (
	"testing"
	"time"
)

type Period struct {
	StartTime, EndTime time.Time
}

func (p Period) Duration() time.Duration {
	return p.EndTime.Sub(p.StartTime)
}

func TestGetDuration(t *testing.T) {
	now := time.Now()
	tests := map[string]struct {
		period Period
		want   time.Duration
	}{
		"test1": {
			period: Period{
				StartTime: now,
				EndTime:   now.Add(time.Hour),
			},
			want: time.Hour,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.period.Duration()
			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want) // want `Prefer "test.period.Duration\(\) = %v, want %v" format for this failure message`
			}
		})
	}
}
