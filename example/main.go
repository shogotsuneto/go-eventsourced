package main

import (
	"fmt"
	"log"

	"github.com/shogotsuneto/go-eventsourced"
	"github.com/shogotsuneto/go-eventsourced/locked"
)

// UserState represents our application state
type UserState struct {
	Name  string
	Age   int
	Email string
}

// Apply implements eventsourced.State[UserState] - mutates the receiver
func (u *UserState) Apply(event eventsourced.Event) error {
	switch e := event.(type) {
	case UserNameChanged:
		u.Name = e.NewName
	case UserAgeChanged:
		u.Age = e.NewAge
	case UserEmailChanged:
		u.Email = e.NewEmail
	default:
		return fmt.Errorf("unknown event type: %s", event.Type())
	}
	return nil
}

// Clone implements eventsourced.State[*UserState] - creates a deep copy
func (u *UserState) Clone() *UserState {
	return &UserState{
		Name:  u.Name,
		Age:   u.Age,
		Email: u.Email,
	}
}

// UserEvent implementations
type UserNameChanged struct {
	NewName string
}

func (e UserNameChanged) Type() string { return "UserNameChanged" }

type UserAgeChanged struct {
	NewAge int
}

func (e UserAgeChanged) Type() string { return "UserAgeChanged" }

type UserEmailChanged struct {
	NewEmail string
}

func (e UserEmailChanged) Type() string { return "UserEmailChanged" }

func main() {
	// Initialize the event sourced state with a zero value
	initialState := &UserState{Name: "Unknown", Age: 0, Email: ""}
	
	// The new design is much simpler - no function parameters needed
	es := locked.New(initialState)

	fmt.Println("=== Event Sourced User State Example ===")
	fmt.Printf("Initial state: %+v\n", es.GetState())

	// Apply some events
	events := []eventsourced.Event{
		UserNameChanged{NewName: "Alice"},
		UserAgeChanged{NewAge: 30},
		UserEmailChanged{NewEmail: "alice@example.com"},
	}

	for _, event := range events {
		if err := es.Apply(event); err != nil {
			log.Fatalf("Error applying event %s: %v", event.Type(), err)
		}
		fmt.Printf("After %s: %+v\n", event.Type(), es.GetState())
	}

	// Demonstrate state cloning to avoid mutation
	fmt.Println("\n=== State Cloning Example ===")
	currentState := es.GetState()
	// GetState now automatically clones the state using the state's Clone method
	clonedState := es.GetState()
	
	fmt.Printf("First call to GetState: %+v\n", currentState)
	fmt.Printf("Second call to GetState: %+v\n", clonedState)
	
	// Modify one of the returned states to show they're independent
	currentState.Name = "Modified State"
	fmt.Printf("After modifying first returned state - First: %+v\n", currentState)
	fmt.Printf("After modifying first returned state - Second: %+v\n", clonedState)
	fmt.Printf("ES GetState() returns fresh clone: %+v\n", es.GetState())
}