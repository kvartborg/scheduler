package scheduler_test

import (
	"testing"
	"time"

	"github.com/kvartborg/scheduler"
)

func TestScheduler(t *testing.T) {
	s := scheduler.New()

	cases := []*struct {
		name     string
		duration time.Duration
		ticks    int
	}{
		{"every millisecond", time.Millisecond, 0},
		{"every 100 milliseconds", 100 * time.Millisecond, 0},
		{"every second", time.Second, 0},
		{"every 10 seconds", 10 * time.Second, 0},
		{"every hour", time.Hour, 0},
	}

	for _, c := range cases {
		r := c
		s.Every(r.duration, func(time.Duration) {
			r.ticks++
		})
	}

	s.Evaluate(time.Hour)

	for _, c := range cases {
		expected := int(time.Hour / c.duration)
		if expected != c.ticks {
			t.Errorf(
				"%s did not get invoked at correct tick rate, expected %d ticks; got %d ticks",
				c.name,
				expected,
				c.ticks,
			)
		}
	}
}
