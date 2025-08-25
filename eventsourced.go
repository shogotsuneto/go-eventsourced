package eventsourced

// Event is plain data; concrete types live in your domain.
type Event interface{ Type() string }

// State constraint for types that can be used with LockedES.
// The state type must know how to apply events and clone itself.
type State[S any] interface {
	Apply(Event) error // mutate receiver
	Clone() S          // deep-enough copy for safe external use
}

// EventApplier mutates state in response to an event.
type EventApplier interface {
	Apply(e Event) error
}

// StateGetter exposes the current state (read-only from the caller's POV).
type StateGetter[S any] interface {
	GetState() S
}

type ES[S any] interface {
	EventApplier
	StateGetter[S]
}