// Package main demonstrates a classic deadlock caused by inconsistent lock
// ordering: two goroutines each lock a pair of mutexes in opposite order,
// so they can end up each holding the lock the other one needs.
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var a, b value
	wg := &sync.WaitGroup{}

	wg.Add(2)
	// goroutine 1 locks a then b ...
	go printSum(&a, &b, wg)
	// ... while goroutine 2 locks b then a. This reversed ordering is what
	// creates the circular wait: each goroutine can grab its first lock
	// (since they're different mutexes) but then blocks forever waiting on
	// the second, which the other goroutine is holding.
	go printSum(&b, &a, wg)
	// Wait() blocks forever here because neither printSum call can ever
	// call wg.Done() - the program deadlocks and the Go runtime will
	// eventually detect and report "all goroutines are asleep - deadlock!".
	wg.Wait()
}

// value pairs an int with its own mutex so it can be safely read/written
// from multiple goroutines.
type value struct {
	mu    sync.Mutex
	value int
}

// printSum locks v1, sleeps (to make the race window wide enough to
// reliably trigger the deadlock), then locks v2 and prints the sum.
// Because callers pass a and b in swapped order, this function is the
// source of the inconsistent lock ordering that deadlocks the program.
func printSum(v1, v2 *value, wg *sync.WaitGroup) {
	defer wg.Done()

	v1.mu.Lock()
	defer v1.mu.Unlock()

	// Simulate work while holding v1's lock, giving the other goroutine
	// time to acquire its own first lock before we try to acquire v2's.
	time.Sleep(2 * time.Second)

	v2.mu.Lock()
	defer v2.mu.Unlock()

	fmt.Printf("sum=%v\n", v1.value+v2.value)
}
