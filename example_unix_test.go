// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix

package signals_test

import (
	"fmt"
	"os"
	"time"

	"github.com/weiwenchen2022/signals"
)

func ExampleNotify() {
	// Set up function on which to receive signal notifications.
	done := make(chan struct{})
	received := 3
	stop := signals.Notify(func(sig os.Signal) bool {
		fmt.Println("Got signal:", sig)

		received--
		if received == 0 {
			close(done)
			return false
		}
		return true
	}, os.Interrupt)
	defer stop()

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		panic(err)
	}

	// On a Unix-like system, pressing Ctrl+C on a keyboard sends a
	// SIGINT signal to the process of the program in execution.
	//
	// This example simulates that by sending a SIGINT signal to itself.
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		send := 3
		for range ticker.C {
			if err := p.Signal(os.Interrupt); err != nil {
				panic(err)
			}

			send--
			if send == 0 {
				return
			}
		}
	}()

	select {
	case <-time.After(300*time.Millisecond + 10*time.Millisecond):
		panic("missed signal")
	case <-done:
		stop() // stop receiving signal notifications as soon as possible.
	}

	// Output:
	// Got signal: interrupt
	// Got signal: interrupt
	// Got signal: interrupt
}
