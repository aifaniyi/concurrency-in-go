# concurrency-in-go

A collection of small, self-contained Go programs exploring core
concurrency primitives and patterns: goroutines, channels, `sync.Mutex`,
`sync.WaitGroup`, `sync.Cond`, and pipelines — along with classic failure
modes like deadlock and starvation.

Each example lives in its own `main` package and can be run independently.

## Requirements

- Go 1.12 or later

## Project structure

| Path | Demonstrates |
| --- | --- |
| [ch01/deadlock](ch01/deadlock/main.go) | A deadlock caused by two goroutines acquiring the same two mutexes in opposite order |
| [ch01/starvation](ch01/starvation/main.go) | Goroutine starvation when one worker holds a shared mutex far more efficiently than another |
| [ch03/channels](ch03/channels/main.go) | A producer/consumer pipeline with graceful shutdown on `SIGINT`/`SIGTERM` |
| [ch03/waitgroup](ch03/waitgroup/main.go) | Basic `sync.WaitGroup` usage to wait for a fixed pool of goroutines |
| [ch04/pipelines](ch04/pipelines/main.go) | Chaining concurrent pipeline stages together via channels, with cancellation support |
| [cmd/bank_transaction](cmd/bank_transaction/main.go) | Protecting a shared invariant (balance total) across concurrent updates with a mutex |
| [cmd/channel-wait](cmd/channel-wait/main.go) | Using an unbuffered channel as a manual join/completion signal |
| [cmd/channels](cmd/channels/main.go) | The blocking, synchronous nature of unbuffered channel sends/receives |
| [cmd/condition_variable](cmd/condition_variable/main.go) | `sync.Cond` for waking a goroutine once a shared condition becomes true |
| [cmd/periodic](cmd/periodic/main.go) | A periodic background goroutine that shuts down cooperatively via a shared flag |
| [cmd/shared_counter](cmd/shared_counter/main.go) | Guarding a shared counter from a data race with a mutex |
| [cmd/unbuffered-deadlock](cmd/unbuffered-deadlock/main.go) | The simplest possible deadlock: sending on an unbuffered channel with no receiver |
| [cmd/waitgroup](cmd/waitgroup/main.go) | Fanning out simulated RPC calls and joining on them with a `WaitGroup` |
| [utils.go](utils.go) | Shared helper (`TimeConsumingOperation`) used by the `cmd/` examples to simulate variable-latency work |

## Running an example

Each example is a separate `main` package, so run it directly with `go run`:

```sh
go run ./ch01/deadlock
go run ./cmd/condition_variable
```

Some examples (like `ch01/deadlock` and `cmd/unbuffered-deadlock`) are
*intentionally* broken — they will hang or crash with a `fatal error: all
goroutines are asleep - deadlock!` to illustrate the failure mode.
