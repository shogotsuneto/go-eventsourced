package eventsourced

import "sync"

// LockedES is a minimal ES[S] implementation with internal locking.
type LockedES[S any] struct {
	mu    sync.RWMutex
	state S
	apply func(*S, Event) error
}

func New[S any](zero S, apply func(*S, Event) error) *LockedES[S] {
	return &LockedES[S]{state: zero, apply: apply}
}

func (x *LockedES[S]) Apply(e Event) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.apply(&x.state, e)
}

func (x *LockedES[S]) GetState() S {
	x.mu.RLock()
	defer x.mu.RUnlock()
	return x.state
}