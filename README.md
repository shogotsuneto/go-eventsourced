# go-eventsourced
Minimal in-memory event-sourced state management for Go with a clean State interface constraint design.

## Usage

State types implement the `eventsourced.State[S]` interface to define their own event handling and cloning behavior:

```go
package main

import (
    "fmt"
    "github.com/shogotsuneto/go-eventsourced"
    "github.com/shogotsuneto/go-eventsourced/locked"
)

// State type implements the State constraint interface
type CounterState struct {
    Value int
}

// Apply mutates the receiver based on the event
func (s *CounterState) Apply(event eventsourced.Event) error {
    switch e := event.(type) {
    case IncrementEvent:
        s.Value += e.Amount
    default:
        return fmt.Errorf("unknown event: %s", event.Type())
    }
    return nil
}

// Clone creates a deep copy for safe external use
func (s *CounterState) Clone() *CounterState {
    return &CounterState{Value: s.Value}
}

// Event implementation
type IncrementEvent struct {
    Amount int
}

func (e IncrementEvent) Type() string { return "increment" }

func main() {
    // Clean constructor - only takes initial state
    es := locked.New(&CounterState{Value: 0})
    
    es.Apply(IncrementEvent{Amount: 5})
    fmt.Printf("State: %+v\n", es.GetState()) // State: &{Value:5}
}
```

## Key Interfaces

The library provides a type-safe foundation through interface constraints:

```go
// Events are plain data with type identification
type Event interface{ Type() string }

// State constraint - types must implement event handling and cloning
type State[S any] interface {
    Apply(Event) error // mutate receiver
    Clone() S          // deep-enough copy for safe external use
}

// Complete event sourcing interface
type ES[S any] interface {
    Apply(Event) error
    GetState() S
}
```

## Features

- **Type-Safe**: Generic constraints ensure required methods exist at compile time
- **Clean API**: Constructor only needs initial state, no function parameters
- **Thread-Safe**: Read-write locks for optimal concurrent access
- **Better Encapsulation**: State types own their behavior (applying events and cloning)
- **Go-Idiomatic**: Methods on types rather than function passing
- **Minimal Dependencies**: Standard library only
- **In-Memory First**: Focused on lightweight, embedded use without external persistence

## Architecture

- **Core Interfaces** (`eventsourced.go`): Event sourcing contracts and type constraints
- **Thread-Safe Implementation** (`locked/`): Generic `LockedES[S State[S]]` with mutex protection
- **Clean Separation**: Interfaces in root package, implementations in sub-packages

## Example

See the [example directory](example/) for a complete working example with multiple event types and state cloning demonstration.
