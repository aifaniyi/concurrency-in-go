// Package main demonstrates using an unbuffered channel as a completion
// signal - a manual alternative to sync.WaitGroup. Each of 10 goroutines
// sends a value on the shared "done" channel when it finishes; main
// receives exactly 10 times to know all of them are done. Because each
// worker sleeps for a different duration (x seconds), the receives in main
// arrive out of program-launch order but always in ascending completion
// order (worker 0 first, worker 9 last).
package main

import (
	"fmt"
	"time"
)

func main() {
	// Unbuffered: each send in the goroutines below blocks until main
	// performs a matching receive, which is fine here since main is always
	// ready to receive in its own loop.
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(x int, end chan bool) {
			// do work
			time.Sleep(time.Duration(x) * time.Second)
			fmt.Printf("worker %d done after %d seconds\n", x, x)

			// send result
			end <- true
		}(i, done)
	}

	// Receiving exactly 10 times - one per launched goroutine - is what
	// makes this a correct join; receiving fewer times would let main exit
	// (and the process terminate) before all workers finish.
	for i := 0; i < 10; i++ {
		<-done
	}
	println("done")
}
