// Package main demonstrates a periodic background goroutine that shuts
// itself down cooperatively when it observes a shared "done" flag, rather
// than being killed abruptly. The flag is a plain bool shared between
// goroutines, so every read and write goes through mu to avoid a data race
// (a bare `var done bool` read/written concurrently without synchronization
// would be undefined behavior under the Go memory model, even though the
// bug wouldn't necessarily show up every run).
package main

import (
	"sync"
	"time"
)

var (
	done bool
	mu   sync.Mutex
)

func main() {
	println("started")
	go periodic()

	// Let the periodic goroutine tick for a while before requesting shutdown.
	time.Sleep(5 * time.Second)
	mu.Lock()
	println("initiating shutdown")
	done = true
	mu.Unlock()

	// Give periodic() time to notice `done` on its next tick and exit
	// before main itself returns (main returning ends the whole program,
	// killing any goroutines still running).
	time.Sleep(3 * time.Second)
}

// periodic prints "tick" once a second forever, checking after each sleep
// whether shutdown has been requested. This poll-on-a-timer approach is
// simple but means shutdown can be observed up to ~1 second late; a channel
// close would notify it immediately instead.
func periodic() {
	for {
		println("tick")

		time.Sleep(1 * time.Second)
		mu.Lock()
		if done {
			println("shutting down...")
			mu.Unlock()
			return
		}
		mu.Unlock()
	}
}
