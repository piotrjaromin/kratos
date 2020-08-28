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
	pacerSwitched  bool
	start          time.Time
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
		pacerSwitched:  false,
		start:          time.Now(),
	}
}

func (p *rampUpPacer) Pace(elapsed time.Duration, hits uint64) (time.Duration, bool) {
	if elapsed > p.rampUpDuration && !p.pacerSwitched {
		fmt.Println("Pacer switch")
		p.currentPacer = p.constPacer
		p.pacerSwitched = true

		e := time.Now().Sub(p.start)
		fmt.Printf("Time taken %+v\n", e)
	}

	t, b := p.currentPacer.Pace(elapsed, hits)
	fmt.Printf("!!! Pace is %+v and %+v, paccer used %+v\n", t, b, p.currentPacer)
	return t, b
	// return p.currentPacer.Pace(elapsed, hits)
}

func (p *rampUpPacer) Rate(elapsed time.Duration) float64 {
	fmt.Printf("--- Got rate %+v\n", p.currentPacer.Rate(elapsed))
	return p.currentPacer.Rate(elapsed)
}
