package scheduler

import "time"

type options struct {
	handlers map[time.Duration][]Handler
}

type Option func(options *options)
