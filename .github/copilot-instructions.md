# go-eventsourced

Minimal in-memory event-sourced state management library for Go with clean State interface constraint design. This is a pure Go library with zero external dependencies that provides type-safe event sourcing through interface constraints.

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### Prerequisites
- Requires Go 1.24.6 or later (specified in go.mod)
- No external dependencies - uses standard library only
- No special tooling required beyond standard Go tools

### Bootstrap and Build Process
- `go mod tidy` -- takes 7 seconds. Downloads any missing standard library dependencies.
- `go build ./...` -- takes 6 seconds. NEVER CANCEL. Builds all packages including root, locked/, and example/.
- `go mod download` -- takes <1 second (no external dependencies to download)

### Testing
- `go test ./...` -- takes 4 seconds. NEVER CANCEL. Runs all tests (only locked/ has tests).
- `go test -v ./locked` -- takes <1 second. Run tests with verbose output for detailed validation.
- `go test -race ./...` -- takes 10 seconds. NEVER CANCEL. Tests with race condition detection (important for concurrent access patterns).

### Code Quality
- `go vet ./...` -- takes <1 second. Always run to catch common Go issues.
- `go fmt ./...` -- takes <1 second. Always run to ensure consistent formatting.
- No external linting tools configured (golangci-lint not present)

### Run the Example
- `cd example && go run main.go` -- takes <1 second. Demonstrates event sourcing with user state management.
- Expected output shows initial state, event applications, and state cloning behavior.

## Validation

### Always Run This Complete Validation After Making Changes
1. Build: `go build ./...` (must succeed)
2. Test: `go test ./...` (all tests must pass)
3. Race test: `go test -race ./...` (must pass - critical for concurrent access)
4. Code quality: `go vet ./...` && `go fmt ./...`
5. Example: `cd example && go run main.go` (must run without errors)

### Manual Validation Scenarios
After making any changes to the core library, ALWAYS verify these scenarios work:
1. **State Creation**: Create new event-sourced state with initial values
2. **Event Application**: Apply multiple different event types and verify state changes
3. **State Cloning**: Verify GetState() returns independent copies that can be modified without affecting internal state
4. **Concurrent Access**: Verify multiple goroutines can safely read/write simultaneously (race detector must pass)

### Test Suite Validation
The comprehensive test suite in `locked/locked_test.go` validates all core functionality:
- **State creation**: `TestNew` verifies initial state setup
- **Event application**: `TestApply` and `TestApplyError` test event processing and error handling
- **State cloning**: `TestStateCloning` ensures GetState() returns independent copies
- **Concurrent access**: `TestConcurrentAccess` validates thread safety with goroutines

Run `go test -v ./locked` to see detailed test output and verify all scenarios pass.

## Key Architecture Components

### Core Files
- **`eventsourced.go`**: Core interfaces (Event, State[S], ES[S]) - type constraints for event sourcing
- **`locked/locked.go`**: Thread-safe LockedES[S] implementation with mutex protection  
- **`locked/locked_test.go`**: Comprehensive tests including concurrent access and state cloning
- **`example/main.go`**: Complete working example with UserState demonstrating all features

### Key Interfaces to Understand
```go
// Events must implement Type() for identification
type Event interface{ Type() string }

// State types must handle events and provide cloning
type State[S any] interface {
    Apply(Event) error // mutate receiver based on event
    Clone() S          // create deep copy for safe external use
}

// Complete event sourcing interface
type ES[S any] interface {
    Apply(Event) error // apply event to internal state
    GetState() S       // return cloned state for safe external access
}
```

### Threading and Concurrency
- **CRITICAL**: LockedES uses sync.RWMutex for thread safety
- Apply() operations use write locks (exclusive)
- GetState() operations use read locks (shared)
- State cloning prevents external mutation of internal state
- Always run race detector when testing concurrent access patterns

## Common Development Tasks

### Adding New Event Types
1. Define event struct implementing `Event` interface with `Type() string` method
2. Update state's `Apply()` method to handle the new event type in switch statement
3. Test event application and state changes
4. ALWAYS test with race detector if the event will be used concurrently

### Adding New State Types
1. Define state struct with required fields
2. Implement `Apply(Event) error` method handling all relevant events
3. Implement `Clone() S` method creating deep copy of all fields
4. Create LockedES instance with `locked.New(initialState)`
5. Write tests covering event application, state cloning, and concurrent access

### Testing Changes
- Unit tests go in `locked/locked_test.go` - follow existing pattern
- Test structure: TestState with Counter/Message fields, TestEvent with EventType/Value
- ALWAYS test concurrent access with goroutines
- ALWAYS test state cloning independence
- ALWAYS run with race detector

## Repository Structure
```
.
├── eventsourced.go        # Core interfaces and type constraints
├── locked/                # Thread-safe implementation package
│   ├── locked.go         # LockedES[S] with mutex protection
│   └── locked_test.go    # Comprehensive test suite
├── example/              # Working example directory
│   └── main.go          # UserState example with multiple event types
├── go.mod               # Go 1.24.6, no external dependencies
├── README.md            # API documentation and usage examples
└── LICENSE              # MIT license
```

## Common Command Reference

### Development Workflow
```bash
# Full validation sequence (run this before committing)
go mod tidy && go build ./... && go test ./... && go test -race ./... && go vet ./... && go fmt ./...

# Quick test
go test ./...

# Detailed test output  
go test -v ./locked

# Run example
cd example && go run main.go

# Check for race conditions (CRITICAL for concurrent code)
go test -race ./...
```

### Timing Expectations
- `go mod tidy`: ~7 seconds
- `go build ./...`: ~6 seconds  
- `go test ./...`: ~4 seconds
- `go test -race ./...`: ~10 seconds
- `go vet ./...`: <1 second
- `go fmt ./...`: <1 second
- Example run: <1 second

**NEVER CANCEL** any of these commands - they complete quickly but race testing is critical for correctness.