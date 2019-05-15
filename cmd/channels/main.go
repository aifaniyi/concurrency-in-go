// Package main demonstrates that sending on an unbuffered channel blocks
// until another goroutine is ready to receive. The receiving goroutine
// deliberately sleeps for 1 second before it calls <-c, so the send in
// main is forced to block for roughly that long - the printed duration
// makes the blocking behavior of unbuffered channels directly observable.
package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan bool)
	go func() {
		// Delay the receive so main's send below has something to visibly
		// wait on.
		time.Sleep(1 * time.Second)
		<-c
	}()

	start := time.Now()
	// This send blocks until the goroutine above reaches its receive,
	// roughly 1 second later - unbuffered channel sends and receives are
	// synchronous rendezvous points, not asynchronous mailboxes.
	c <- true
	fmt.Printf("send took %v\n", time.Since(start))
}
