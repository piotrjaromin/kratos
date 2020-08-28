package attack

import (
	"fmt"
	"io"
	"os"
	"os/signal"
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

	reporter, report, err := getReport(typ, bucketsStr)
	if err != nil {
		return err
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
			if err = writeReport(reporter, rc, out); err != nil {
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

	img, ok := report.(*ImgReport)

	if ok {
		return img.Draw()
	}

	return writeReport(reporter, rc, out)
}

func decoder(contets io.Reader) (vegeta.Decoder, error) {
	dec := vegeta.DecoderFor(contets)
	if dec == nil {
		return nil, fmt.Errorf("encode: can't detect encoding for report string")
	}

	return vegeta.NewRoundRobinDecoder(dec), nil
}

func getReport(typ, bucketsStr string) (vegeta.Reporter, vegeta.Report, error) {
	switch typ {

	case "text":
		var m vegeta.Metrics
		return vegeta.NewTextReporter(&m), &m, nil

	case "json":
		var m vegeta.Metrics
		return vegeta.NewJSONReporter(&m), &m, nil

	case "img":
		img, m := NewImgReport()
		return vegeta.NewJSONReporter(m), img, nil

	case "hdrplot":
		var m vegeta.Metrics
		return vegeta.NewHDRHistogramPlotReporter(&m), &m, nil

	case "hist":
		var hist vegeta.Histogram
		if err := hist.Buckets.UnmarshalText([]byte(bucketsStr)); err != nil {
			return nil, nil, err
		}
		return vegeta.NewHistogramReporter(&hist), &hist, nil

	default:
		return nil, nil, fmt.Errorf("unknown report type: %q", typ)
	}

}

func writeReport(r vegeta.Reporter, rc vegeta.Closer, out io.Writer) error {
	if rc != nil {
		rc.Close()
	}
	return r.Report(out)
}
