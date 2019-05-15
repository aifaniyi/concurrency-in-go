// Package main demonstrates protecting a shared invariant (that the total
// of two "account balances" never changes) with a single mutex shared by
// every goroutine that touches either balance. Two goroutines transfer
// money back and forth between alice and bob, and a checker loop on main
// continuously verifies the invariant holds. Because every read and write
// of alice/bob happens under the same mu, the checker will never observe a
// torn/inconsistent state.
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex

	alice := 10000
	bob := 10000
	total := alice + bob

	// Transfers $1 from alice to bob, 1000 times.
	go func() {
		for i := 0; i < 1000; i++ {
			mu.Lock()
			alice -= 1
			bob += 1
			mu.Unlock()
		}
	}()

	// Concurrently transfers $1 from bob to alice, 1000 times. Locking the
	// same mu as the goroutine above ensures the two transfers never
	// interleave mid-update (e.g. one goroutine reading alice while the
	// other is only halfway through updating it).
	go func() {
		for i := 0; i < 1000; i++ {
			mu.Lock()
			bob -= 1
			alice += 1
			mu.Unlock()
		}
	}()

	// Repeatedly checks, for up to 1 second, that alice+bob still equals
	// the original total. Taking the lock here is essential: without it,
	// this goroutine could read alice and bob at two different points in
	// time relative to the transfer goroutines and see a false violation
	// even though each individual transfer was atomic.
	start := time.Now()
	for time.Since(start) < 1*time.Second {
		mu.Lock()
		if alice+bob != total {
			fmt.Printf("violation: alice=%d bob=%d total=%d", alice, bob, total)
		}
		mu.Unlock()
	}
}
