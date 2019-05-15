// Package main demonstrates sync.Cond as a way for the main goroutine to
// sleep efficiently until a condition becomes true, instead of busy-polling
// a shared variable. It simulates a 10-voter election: main waits until
// either 5 "yes" votes have come in (an early majority) or all 10 voters
// have finished, whichever happens first, then reports the outcome.
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/aifaniyi/concurrency-in-go"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	count := 0    // number of "yes" votes so far
	finished := 0 // number of voters that have finished voting
	var mu sync.Mutex
	// cond pairs a Locker with a wait/wake mechanism: goroutines call
	// cond.Wait() to sleep until cond.Broadcast()/Signal() wakes them, all
	// while cond's associated mutex protects the shared state (count,
	// finished) being waited on.
	cond := sync.NewCond(&mu)

	for i := 0; i < 10; i++ {
		go func() {
			vote := requestVote()

			// lock then modify conditions that can change outcome i.e count and finished
			mu.Lock()
			defer mu.Unlock()
			if vote {
				count++
			}
			finished++
			cond.Broadcast() // wake up people waiting on the condition ; wakes all waiters  (signal wakes only one waiter)
		}()
	}

	// get lock check condition and wait if false or continue if true
	mu.Lock()
	// Looping (not "if") around Wait() is required: Broadcast can wake this
	// goroutine spuriously or when the condition still isn't true yet (e.g.
	// after a vote that doesn't push count to 5), so the condition must be
	// re-checked every time Wait returns.
	for count < 5 && finished != 10 {
		cond.Wait() // if false relesaes lock mu puts us to sleep to be later woken up by the broadcast
	}
	if count >= 5 {
		fmt.Printf("won election with %d votes\n", count)
	} else {
		fmt.Printf("lost election with %d votes\n", count)
	}
	mu.Unlock()
}

// requestVote simulates the latency of collecting one vote (a random delay
// via TimeConsumingOperation) and always returns true - every simulated
// voter votes "yes".
func requestVote() bool {
	fmt.Println("voting")
	_ = concurrency.TimeConsumingOperation()
	fmt.Println("voted")
	return true
}
