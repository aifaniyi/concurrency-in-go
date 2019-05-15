// Package main demonstrates the simplest possible deadlock: sending on an
// unbuffered channel from the only goroutine in the program (main itself).
// An unbuffered send blocks until some other goroutine performs the
// matching receive, but here there is no other goroutine - the receive on
// the next line can never run because main is permanently blocked on the
// line before it. The Go runtime detects this ("all goroutines are asleep -
// deadlock!") and terminates the program.
package main

func main() {
	c := make(chan bool)
	c <- true // blocks waiting for receiver and never returns
	<-c
}
