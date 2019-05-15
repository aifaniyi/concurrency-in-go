// Package concurrency provides small shared helpers used by the example
// programs in cmd/ to simulate realistic, variable-latency work (like a
// network call or disk I/O) without needing an actual external dependency.
package concurrency

import (
	"math/rand"
	"time"
)

// init seeds the global math/rand source with the current time so that
// TimeConsumingOperation produces a different sequence of durations on each
// run of the program, rather than the same deterministic sequence every
// time (math/rand's default seed is fixed).
func init() {
	rand.Seed(time.Now().UnixNano())
}

// TimeConsumingOperation blocks the calling goroutine for a random duration
// between 100 and 1000 milliseconds (inclusive) and returns that duration
// in milliseconds. It stands in for any real-world operation whose latency
// varies unpredictably, letting the examples demonstrate concurrency
// patterns without depending on actual I/O.
func TimeConsumingOperation() int {
	min := 100
	max := 1000
	val := rand.Intn(max-min+1) + min
	time.Sleep(time.Duration(val) * time.Millisecond)
	return val
}
