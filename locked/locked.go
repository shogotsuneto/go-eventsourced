package locked

import (
	"sync"

	"github.com/shogotsuneto/go-eventsourced"
)

// LockedES is a minimal ES[S] implementation with internal locking.
// S must be a pointer type that implements eventsourced.State[S].
type LockedES[S eventsourced.State[S]] struct {
	mu    sync.RWMutex
	state S
}

func New[S eventsourced.State[S]](zero S) *LockedES[S] {
	return &LockedES[S]{state: zero}
}

func (x *LockedES[S]) Apply(e eventsourced.Event) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.state.Apply(e)
}

func (x *LockedES[S]) GetState() S {
	x.mu.RLock()
	defer x.mu.RUnlock()
	return x.state.Clone()
}