# go-eventsourced
Minimal in-memory event-sourced state management for Go.

## Usage

```go
package main

import (
    "fmt"
    "github.com/shogotsuneto/go-eventsourced"
)

type MyState struct {
    Counter int
}

type IncrementEvent struct {
    Amount int
}

func (e IncrementEvent) Type() string { return "increment" }

func applyEvent(state *MyState, event eventsourced.Event) error {
    switch e := event.(type) {
    case IncrementEvent:
        state.Counter += e.Amount
    }
    return nil
}

func main() {
    es := eventsourced.New(MyState{Counter: 0}, applyEvent)
    
    es.Apply(IncrementEvent{Amount: 5})
    fmt.Printf("State: %+v\n", es.GetState()) // State: {Counter:5}
}
```

## Features

- Thread-safe state management with read-write locks
- Generic type support for any state type
- Simple event application pattern
- Minimal dependencies (only standard library)

## Example

See the [example directory](example/) for a complete working example with state cloning.
