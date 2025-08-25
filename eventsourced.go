package eventsourced

// Event is plain data; concrete types live in your domain.
type Event interface{ Type() string }

// EventApplier mutates state in response to an event.
type EventApplier interface {
	Apply(e Event) error
}

// StateGetter exposes the current state (read-only from the caller's POV).
// If S is a pointer, you're exposing the same pointer.
// Wrap with your own deep-clone when needed.
type StateGetter[S any] interface {
	GetState() S
}

type ES[S any] interface {
	EventApplier
	StateGetter[S]
}