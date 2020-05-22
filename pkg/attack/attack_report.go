package attack

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/piotrjaromin/kratos/pkg/utils"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func Report(attackData io.Reader, typ, bucketsStr string, every time.Duration) error {
	out, err := utils.File("stdout", true)
	if err != nil {
		return err
	}
	defer out.Close()

	dec, err := decoder(attackData)
	if err != nil {
		return err
	}

	var (
		rep    vegeta.Reporter
		report vegeta.Report
	)

	switch typ {
	case "text":
		var m vegeta.Metrics
		rep, report = vegeta.NewTextReporter(&m), &m
	case "json":
		var m vegeta.Metrics
		if bucketsStr != "" {
			m.Histogram = &vegeta.Histogram{}
			if err := m.Histogram.Buckets.UnmarshalText([]byte(bucketsStr)); err != nil {
				return err
			}
		}
		rep, report = vegeta.NewJSONReporter(&m), &m
	case "hdrplot":
		var m vegeta.Metrics
		rep, report = vegeta.NewHDRHistogramPlotReporter(&m), &m
	default:
		switch {
		case strings.HasPrefix(typ, "hist"):
			var hist vegeta.Histogram
			if bucketsStr == "" { // Old way
				if len(typ) < 6 {
					return fmt.Errorf("bad buckets: '%s'", typ[4:])
				}
				bucketsStr = typ[4:]
			}
			if err := hist.Buckets.UnmarshalText([]byte(bucketsStr)); err != nil {
				return err
			}
			rep, report = vegeta.NewHistogramReporter(&hist), &hist
		default:
			return fmt.Errorf("unknown report type: %q", typ)
		}
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)

	var ticks <-chan time.Time
	if every > 0 {
		ticker := time.NewTicker(every)
		defer ticker.Stop()
		ticks = ticker.C
	}

	rc, _ := report.(vegeta.Closer)
decode:
	for {
		select {
		case <-sigch:
			break decode
		case <-ticks:
			if err = writeReport(rep, rc, out); err != nil {
				return err
			}
		default:
			var r vegeta.Result
			if err = dec.Decode(&r); err != nil {
				if err == io.EOF {
					break decode
				}
				return err
			}

			report.Add(&r)
		}
	}

	return writeReport(rep, rc, out)
}

func decoder(contets io.Reader) (vegeta.Decoder, error) {
	dec := vegeta.DecoderFor(contets)
	if dec == nil {
		return nil, fmt.Errorf("encode: can't detect encoding for report string")
	}

	return vegeta.NewRoundRobinDecoder(dec), nil
}

func writeReport(r vegeta.Reporter, rc vegeta.Closer, out io.Writer) error {
	if rc != nil {
		rc.Close()
	}
	return r.Report(out)
}
