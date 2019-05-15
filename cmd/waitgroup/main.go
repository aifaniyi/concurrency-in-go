// Package main demonstrates fanning out simulated RPC calls across several
// goroutines and using a sync.WaitGroup to block until they've all
// completed, without needing to know how long each one takes in advance.
package main

import (
	"fmt"
	"sync"

	concurrency "github.com/aifaniyi/concurrency-in-go"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		// Registered before the goroutine starts so Wait() can't return
		// prematurely due to a race between Add and the goroutine's Done.
		wg.Add(1)
		go func(x int) {
			sendRPC(x)
			wg.Done()
		}(i)
	}

	// Blocks until all 5 simulated RPCs have completed and called Done().
	wg.Wait()
}

// sendRPC simulates an RPC call with a random, variable-length duration
// (via TimeConsumingOperation) and reports how long it took.
func sendRPC(x int) {
	val := concurrency.TimeConsumingOperation()
	fmt.Printf("routine %d: operation completed in %d ms\n", x, val)
}
