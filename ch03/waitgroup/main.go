// Package main demonstrates the basic sync.WaitGroup pattern: launch a fixed
// number of goroutines, and block the main goroutine until all of them have
// finished. Because each worker sleeps for the same duration and they all
// run concurrently, the total elapsed time printed at the end is roughly
// 2 seconds (the duration of one worker) rather than 10 seconds (5 workers
// run sequentially), proving the work overlaps.
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	wg := &sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		// Add(1) must happen before the goroutine is started, not inside
		// it, otherwise Wait() could race ahead and return before all
		// goroutines have registered themselves.
		wg.Add(1)
		go doWork(wg, 2*time.Second, i)
	}

	// Blocks until every goroutine has called wg.Done().
	wg.Wait()
	fmt.Printf("all routines completed in %v", time.Since(start))
}

// doWork simulates a unit of work of a fixed duration and signals the
// WaitGroup when done, via defer, so Done() still runs even if this
// function were to panic or return early.
func doWork(wg *sync.WaitGroup, workTime time.Duration, id int) {
	defer wg.Done()

	fmt.Printf("starting worker: %d\n", id)
	time.Sleep(workTime)
	fmt.Printf("stopping worker: %d\n", id)
}
