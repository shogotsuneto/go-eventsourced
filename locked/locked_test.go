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

// Test state
type TestState struct {
	Counter int
	Message string
}

// Test apply function
func testApplyFunc(state *TestState, event eventsourced.Event) error {
	switch e := event.(type) {
	case TestEvent:
		if e.EventType == "increment" {
			state.Counter += e.Value
		} else if e.EventType == "set_message" {
			state.Message = "test message"
		} else if e.EventType == "error" {
			return errors.New("test error")
		}
	}
	return nil
}

func TestNew(t *testing.T) {
	initialState := TestState{Counter: 0, Message: ""}
	es := New(initialState, testApplyFunc, nil) // nil clone function should use identity

	if es == nil {
		t.Fatal("New() returned nil")
	}

	state := es.GetState()
	if state.Counter != 0 || state.Message != "" {
		t.Errorf("Initial state incorrect: %+v", state)
	}
}

func TestApply(t *testing.T) {
	initialState := TestState{Counter: 0, Message: ""}
	es := New(initialState, testApplyFunc, nil) // nil clone function should use identity

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
	initialState := TestState{Counter: 0, Message: ""}
	es := New(initialState, testApplyFunc, nil) // nil clone function should use identity

	// Test error case
	event := TestEvent{EventType: "error", Value: 0}
	err := es.Apply(event)
	if err == nil {
		t.Error("Apply() should have returned an error")
	}
}

func TestGetState(t *testing.T) {
	initialState := TestState{Counter: 42, Message: "test"}
	es := New(initialState, testApplyFunc, nil) // nil clone function should use identity

	state := es.GetState()
	if state.Counter != 42 || state.Message != "test" {
		t.Errorf("GetState() returned incorrect state: %+v", state)
	}
}

func TestConcurrentAccess(t *testing.T) {
	initialState := TestState{Counter: 0, Message: ""}
	es := New(initialState, testApplyFunc, nil) // nil clone function should use identity

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

func TestNilApplyFunctionPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when apply function is nil")
		}
	}()

	initialState := TestState{Counter: 0, Message: ""}
	New(initialState, nil, nil) // This should panic
}

func TestCloneFunction(t *testing.T) {
	initialState := TestState{Counter: 42, Message: "test"}
	
	// Test with custom clone function
	cloneFunc := func(s TestState) TestState {
		return TestState{Counter: s.Counter + 100, Message: s.Message + "_cloned"}
	}
	
	es := New(initialState, testApplyFunc, cloneFunc)
	
	state := es.GetState()
	// State should be cloned according to our custom clone function
	if state.Counter != 142 || state.Message != "test_cloned" {
		t.Errorf("Clone function not working correctly: %+v", state)
	}
}

func TestNilCloneFunctionUsesIdentity(t *testing.T) {
	initialState := TestState{Counter: 42, Message: "test"}
	es := New(initialState, testApplyFunc, nil) // nil clone function should use identity

	state := es.GetState()
	// With identity clone function, state should be returned as-is
	if state.Counter != 42 || state.Message != "test" {
		t.Errorf("Identity clone function not working correctly: %+v", state)
	}
}