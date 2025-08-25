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

// Clone creates a deep copy of the UserState
func (u UserState) Clone() UserState {
	return UserState{
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

// applyUserEvent handles state transitions
func applyUserEvent(state *UserState, event eventsourced.Event) error {
	switch e := event.(type) {
	case UserNameChanged:
		state.Name = e.NewName
	case UserAgeChanged:
		state.Age = e.NewAge
	case UserEmailChanged:
		state.Email = e.NewEmail
	default:
		return fmt.Errorf("unknown event type: %s", event.Type())
	}
	return nil
}

func main() {
	// Initialize the event sourced state with a zero value
	initialState := UserState{Name: "Unknown", Age: 0, Email: ""}
	es := locked.New(initialState, applyUserEvent)

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
	clonedState := currentState.Clone()
	
	fmt.Printf("Original state: %+v\n", currentState)
	fmt.Printf("Cloned state: %+v\n", clonedState)
	
	// Modify the clone to show they're independent
	clonedState.Name = "Modified Clone"
	fmt.Printf("After modifying clone - Original: %+v\n", currentState)
	fmt.Printf("After modifying clone - Clone: %+v\n", clonedState)
	fmt.Printf("ES state unchanged: %+v\n", es.GetState())
}