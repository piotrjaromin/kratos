package attack

import (
	"fmt"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type rampUpPacer struct {
	linearPacer vegeta.LinearPacer
	constPacer  vegeta.ConstantPacer
	// pacer which is currently being used
	currentPacer   vegeta.Pacer
	rampUpDuration time.Duration
}

func NewRampUpPacer(rampUpDuration time.Duration, maxRampUpRps int) *rampUpPacer {
	slope := (float64(maxRampUpRps) / float64(rampUpDuration)) * float64(time.Second)
	linearPacer := vegeta.LinearPacer{
		StartAt: vegeta.ConstantPacer{
			Freq: 0,
			Per:  time.Second,
		},
		Slope: slope,
	}

	constPacer := vegeta.ConstantPacer{
		Freq: maxRampUpRps,
		Per:  time.Second,
	}
	return &rampUpPacer{
		linearPacer:    linearPacer,
		constPacer:     constPacer,
		currentPacer:   linearPacer,
		rampUpDuration: rampUpDuration,
	}
}

func (p *rampUpPacer) Pace(elapsed time.Duration, hits uint64) (time.Duration, bool) {
	if elapsed > p.rampUpDuration {
		p.currentPacer = p.constPacer
	}

	return p.currentPacer.Pace(elapsed, hits)
}

func (p *rampUpPacer) Rate(elapsed time.Duration) float64 {
	fmt.Printf("Got rate %+v\n", p.currentPacer.Rate(elapsed))
	return p.currentPacer.Rate(elapsed)
}
