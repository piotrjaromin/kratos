package attack

import (
	"io"
	"os"
	"os/signal"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// Run starts new load test
func Run(out io.Writer, atk *vegeta.Attacker, targeter vegeta.Targeter, opts Opts) error {
	res := atk.Attack(targeter, opts.Rate, opts.Duration, opts.Name)
	enc := vegeta.NewEncoder(out)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	for {
		select {
		case <-sig:
			atk.Stop()
			return nil
		case r, ok := <-res:
			if !ok {
				return nil
			}
			if err := enc.Encode(r); err != nil {
				return err
			}
		}
	}

	return nil
}
