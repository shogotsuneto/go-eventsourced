package locked

import (
	"sync"

	"github.com/shogotsuneto/go-eventsourced"
)

// LockedES is a minimal ES[S] implementation with internal locking.
type LockedES[S any] struct {
	mu    sync.RWMutex
	state S
	apply func(*S, eventsourced.Event) error
	clone func(S) S
}

func New[S any](zero S, apply func(*S, eventsourced.Event) error, clone func(S) S) *LockedES[S] {
	if apply == nil {
		panic("apply function cannot be nil")
	}
	
	if clone == nil {
		// Default to identity function - return state as-is
		clone = func(s S) S { return s }
	}
	
	return &LockedES[S]{state: zero, apply: apply, clone: clone}
}

func (x *LockedES[S]) Apply(e eventsourced.Event) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.apply(&x.state, e)
}

func (x *LockedES[S]) GetState() S {
	x.mu.RLock()
	defer x.mu.RUnlock()
	return x.clone(x.state)
}