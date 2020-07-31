package attack

import (
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// Opts define details of how attack should be performed
type Opts struct {
	Name           string
	HTTP2          bool
	Lazy           bool
	Duration       time.Duration
	Timeout        time.Duration
	Rate           vegeta.Pacer
	Workers        uint64
	MaxWorkers     uint64
	Connections    int
	MaxConnections int
	MaxBody        int64
	Keepalive      bool
}

// DefaultOpts creates default opts for test
func DefaultOpts(name string, duration time.Duration, rate vegeta.Pacer) Opts {
	return Opts{
		Name:           name,
		Duration:       duration,
		Timeout:        time.Second * 3,
		Rate:           rate,
		Workers:        4,
		MaxWorkers:     8,
		Connections:    50,
		MaxConnections: 50,
		MaxBody:        -1,
		Keepalive:      false,
	}
}
