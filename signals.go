// Package signals provides simple wrapper for os/signal to access to incoming signals.
package signals

import (
	"os"
	"os/signal"
	"sync"
)

// Notify delivers incoming signals to f.
// If no signals are provided, all incoming signals will be delivered to f.
// Otherwise, just the provided signals will.
//
// Notify calls f sequentially for incoming signals.
// If f returns false, causes to stop delivering incoming signals to f.
//
// It is allowed to call Notify multiple times with different functions
// and the same signals: each function invoked with copies of incoming
// signals independently.
//
// The stop function stops delivering incoming signals to f.
// It undoes the effect of calls to Notify that returning stop.
//
// The stop function releases resources associated with it, so code should
// call stop as soon as the operations running in this function complete and
// signals no longer need to be delivered to the function.
func Notify(f func(os.Signal) bool, sigs ...os.Signal) (stop func()) {
	c := make(chan os.Signal, len(sigs)+1)
	signal.Notify(c, sigs...)

	var once sync.Once
	stop = func() {
		once.Do(func() {
			signal.Stop(c)
			close(c)
		})
	}

	go func() {
		for s := range c {
			if !f(s) {
				stop()
				return
			}
		}
	}()
	return stop
}
