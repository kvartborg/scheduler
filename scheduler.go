package scheduler

import (
	"time"
)

// Handler is a function receives the duration of which it was invoked by the
// scheduler.
type Handler func(time.Duration)

// Scheduler contains handlers that it will invoke at the correct interval based
// on its the internal clock.
type Scheduler struct {
	paused    bool
	rewind    bool
	clock     time.Duration
	rate      time.Duration
	nextClock map[time.Duration]time.Duration
	handlers  map[time.Duration][]Handler
}

// New create a new instance of a scheduler
func New(opts ...Option) *Scheduler {
	options := &options{
		handlers: map[time.Duration][]Handler{},
	}

	for _, set := range opts {
		set(options)
	}

	scheduler := &Scheduler{
		paused:    true,
		clock:     0,
		rate:      time.Hour,
		nextClock: map[time.Duration]time.Duration{},
		handlers:  options.handlers,
	}

	for rate := range scheduler.handlers {
		if rate < scheduler.rate {
			scheduler.rate = rate
		}

		scheduler.nextClock[rate] = 0
	}

	return scheduler
}

// Clone a scheduler and all its settings
func (scheduler *Scheduler) Clone() *Scheduler {
	return &Scheduler{
		paused:    true,
		clock:     0,
		rate:      scheduler.rate,
		nextClock: scheduler.nextClock,
		handlers:  scheduler.handlers,
	}
}

// Every registers handlers with a given frequency
func (scheduler *Scheduler) Every(frequency time.Duration, handlers ...Handler) {
	if frequency < scheduler.rate {
		scheduler.rate = frequency
	}

	if _, ok := scheduler.handlers[frequency]; !ok {
		scheduler.handlers[frequency] = []Handler{}
		scheduler.nextClock[frequency] = frequency
	}

	scheduler.handlers[frequency] = append(scheduler.handlers[frequency], handlers...)
}

// Cycle runs a single cycle, a cycle will be the lowest frequency registed with
// the Every method
func (scheduler *Scheduler) Cycle(rates ...time.Duration) {
	for rate, handlers := range scheduler.handlers {
		if scheduler.nextClock[rate] > scheduler.clock {
			continue
		}

		for _, handler := range handlers {
			if scheduler.rewind {
				handler(-rate)
			} else {
				handler(rate)
			}
		}

		scheduler.nextClock[rate] += rate
	}

	scheduler.clock += scheduler.rate
}

// Evaluate will run all cycles within a given duration as fast as posible.
func (scheduler *Scheduler) Evaluate(d time.Duration) {
	clock := time.Now()
	var end time.Time

	if d < 0 {
		scheduler.rewind = true
		end = clock.Add(-(d - 1))
	} else {
		scheduler.rewind = false
		end = clock.Add(d + 1)
	}

	for clock.Before(end) {
		scheduler.Cycle()
		clock = clock.Add(scheduler.rate)
	}
}

// Start run the cycles in real time.
func (scheduler *Scheduler) Start() {
	scheduler.paused = false

	for range time.Tick(scheduler.rate) {
		if scheduler.paused {
			continue
		}

		scheduler.Cycle()
	}
}

// Pause the scheduler from running in realtime
func (scheduler *Scheduler) Pause() {
	scheduler.paused = true
}
