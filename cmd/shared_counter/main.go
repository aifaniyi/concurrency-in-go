// Package main demonstrates protecting a shared counter from a data race
// with a mutex. 1000 goroutines each increment the same counter variable;
// without the mutex, the read-modify-write of counter++ is not atomic and
// concurrent increments could be lost. Note this example uses a fixed sleep
// instead of a sync.WaitGroup to wait for the goroutines, which is a race
// in itself (nothing guarantees all 1000 goroutines finish within 1
// second) - see cmd/waitgroup or ch03/waitgroup for the correct pattern.
package main

import (
	"sync"
	"time"
)

func main() {
	counter := 0
	var mu sync.Mutex

	for i := 0; i < 1000; i++ {
		go func() {
			// Lock/Unlock around the increment makes counter++ effectively
			// atomic: only one goroutine can be inside the critical
			// section at a time, so no increments are lost.
			mu.Lock()
			defer mu.Unlock()
			counter++
		}()
	}

	// Best-effort wait for the goroutines to finish; not a real
	// synchronization guarantee.
	time.Sleep(1 * time.Second)
	println(counter)
}
