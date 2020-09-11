package attack

import (
	"time"

	"github.com/piotrjaromin/kratos/pkg/log"
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
	logger         *log.Logger
}

func NewRampUpPacer(logger *log.Logger, rampUpDuration time.Duration, maxRampUpRps int) *rampUpPacer {
	slope := (float64(maxRampUpRps) / float64(rampUpDuration)) * float64(time.Second)
	linearPacer := vegeta.LinearPacer{
		StartAt: vegeta.ConstantPacer{
			Freq: 1, // This cannot be 0, it breaks
			Per:  time.Second,
		},
		Slope: slope,
	}

	constPacer := vegeta.ConstantPacer{
		Freq: maxRampUpRps,
		Per:  time.Second,
	}

	var currentPacer vegeta.Pacer = linearPacer
	if rampUpDuration == time.Duration(0) {
		currentPacer = constPacer
	}

	return &rampUpPacer{
		linearPacer:    linearPacer,
		constPacer:     constPacer,
		currentPacer:   currentPacer,
		rampUpDuration: rampUpDuration,
		pacerSwitched:  false,
		start:          time.Now(),
		logger:         logger,
	}
}

func (p *rampUpPacer) Pace(elapsed time.Duration, hits uint64) (time.Duration, bool) {
	p.logger.Infof("Pace called %+v, %+v", elapsed, hits)
	if elapsed > p.rampUpDuration && !p.pacerSwitched {
		p.logger.Infof("Rampup reached, switching to constant load")
		p.currentPacer = p.constPacer
		p.pacerSwitched = true
	}

	return p.currentPacer.Pace(elapsed, hits)
}

func (p *rampUpPacer) Rate(elapsed time.Duration) float64 {
	p.logger.Infof("Got rate %+v\n", p.currentPacer.Rate(elapsed))
	return p.currentPacer.Rate(elapsed)
}
