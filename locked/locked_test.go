package locked

import (
	"errors"
	"testing"

	"github.com/shogotsuneto/go-eventsourced"
)

// Test Event implementations
type TestEvent struct {
	EventType string
	Value     int
}

func (e TestEvent) Type() string { return e.EventType }

// Test state that implements eventsourced.State[TestState]
type TestState struct {
	Counter int
	Message string
}

// Apply implements eventsourced.State[TestState] - mutates the receiver
func (s *TestState) Apply(event eventsourced.Event) error {
	switch e := event.(type) {
	case TestEvent:
		if e.EventType == "increment" {
			s.Counter += e.Value
		} else if e.EventType == "set_message" {
			s.Message = "test message"
		} else if e.EventType == "error" {
			return errors.New("test error")
		}
	}
	return nil
}

// Clone implements eventsourced.State[*TestState] - returns a deep copy
func (s *TestState) Clone() *TestState {
	return &TestState{
		Counter: s.Counter,
		Message: s.Message,
	}
}

func TestNew(t *testing.T) {
	initialState := &TestState{Counter: 0, Message: ""}
	es := New(initialState)

	if es == nil {
		t.Fatal("New() returned nil")
	}

	state := es.GetState()
	if state.Counter != 0 || state.Message != "" {
		t.Errorf("Initial state incorrect: %+v", state)
	}
}

func TestApply(t *testing.T) {
	initialState := &TestState{Counter: 0, Message: ""}
	es := New(initialState)

	// Test successful event application
	event := TestEvent{EventType: "increment", Value: 5}
	err := es.Apply(event)
	if err != nil {
		t.Errorf("Apply() returned error: %v", err)
	}

	state := es.GetState()
	if state.Counter != 5 {
		t.Errorf("State not updated correctly: expected Counter=5, got %d", state.Counter)
	}
}

func TestApplyError(t *testing.T) {
	initialState := &TestState{Counter: 0, Message: ""}
	es := New(initialState)

	// Test error case
	event := TestEvent{EventType: "error", Value: 0}
	err := es.Apply(event)
	if err == nil {
		t.Error("Apply() should have returned an error")
	}
}

func TestGetState(t *testing.T) {
	initialState := &TestState{Counter: 42, Message: "test"}
	es := New(initialState)

	state := es.GetState()
	if state.Counter != 42 || state.Message != "test" {
		t.Errorf("GetState() returned incorrect state: %+v", state)
	}
}

func TestConcurrentAccess(t *testing.T) {
	initialState := &TestState{Counter: 0, Message: ""}
	es := New(initialState)

	// Test concurrent reads and writes
	done := make(chan bool, 2)

	// Goroutine 1: Apply events
	go func() {
		for i := 0; i < 100; i++ {
			event := TestEvent{EventType: "increment", Value: 1}
			es.Apply(event)
		}
		done <- true
	}()

	// Goroutine 2: Read state
	go func() {
		for i := 0; i < 100; i++ {
			es.GetState()
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	// Final state should have counter = 100
	finalState := es.GetState()
	if finalState.Counter != 100 {
		t.Errorf("Expected final counter to be 100, got %d", finalState.Counter)
	}
}

func TestStateCloning(t *testing.T) {
	initialState := &TestState{Counter: 42, Message: "test"}
	es := New(initialState)
	
	state1 := es.GetState()
	state2 := es.GetState()
	
	// Modify one state to ensure they are independent
	state1.Counter = 999
	state1.Message = "modified"
	
	// state2 should be unaffected
	if state2.Counter != 42 || state2.Message != "test" {
		t.Errorf("States are not properly cloned: state2=%+v", state2)
	}
	
	// Getting state again should return the original, unmodified state
	state3 := es.GetState()
	if state3.Counter != 42 || state3.Message != "test" {
		t.Errorf("Internal state was modified: state3=%+v", state3)
	}
}