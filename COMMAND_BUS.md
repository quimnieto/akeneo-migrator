# Command Bus Architecture

## Overview

The application uses a Command Bus pattern to decouple command execution from command handling. This provides better separation of concerns, testability, and extensibility through middleware.

## Architecture

```
┌─────────────┐
│   CLI       │
│  Commands   │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Command    │
│    Bus      │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Middlewares │ (Logging, etc.)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Command    │
│  Handlers   │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Services   │
└─────────────┘
```

## Components

### 1. Command Bus (`kit/bus/in_memory/command_bus.go`)

The command bus is responsible for:
- Dispatching commands to their handlers
- Executing middleware chain
- Managing handler registration

```go
type CommandBus struct {
    handlers    map[bus.Type]bus.Handler
    middlewares []Middleware
}
```

### 2. Commands

Commands are simple data structures that represent an intention to perform an action:

```go
// Example: SyncAttributeCommand
type SyncAttributeCommand struct {
    Code  string
    Debug bool
}

func (c SyncAttributeCommand) Type() bus.Type {
    return SyncAttributeCommandType
}
```

**Available Commands:**
- `SyncReferenceEntityCommand` - Sync a reference entity
- `SyncProductCommand` - Sync a single product
- `SyncProductHierarchyCommand` - Sync a product hierarchy
- `SyncAttributeCommand` - Sync an attribute

### 3. Command Handlers

Handlers execute the business logic for each command:

```go
type CommandHandler struct {
    service *Service
}

func (h *CommandHandler) Handle(ctx context.Context, msg bus.Message) (bus.Response, error) {
    cmd := msg.(SyncAttributeCommand)
    result, err := h.service.Sync(ctx, cmd.Code)
    return bus.Response{Data: result}, err
}
```

### 4. Middlewares

Middlewares wrap command execution to add cross-cutting concerns:

```go
type Middleware func(ctx context.Context, msg bus.Message, next NextFunc) (bus.Response, error)
```

**Available Middlewares:**
- `Logging` - Logs command execution and timing

## Usage

### Registering Commands

In `bootstrap.go`:

```go
// Create command bus with middlewares
commandBus := inmemory.NewCommandBus(
    middleware.Logging(),
)

// Register handlers
commandBus.Register(
    syncing.SyncReferenceEntityCommandType,
    syncing.NewCommandHandler(referenceEntitySyncer),
)
```

### Dispatching Commands

```go
response, err := app.CommandBus.Dispatch(ctx, attribute_syncing.SyncAttributeCommand{
    Code:  "sku",
    Debug: true,
})

if err != nil {
    return err
}

result := response.Data.(*attribute_syncing.SyncResult)
```

## Benefits

### 1. Separation of Concerns
- CLI layer only knows about commands
- Business logic is isolated in services
- Infrastructure concerns handled by bus

### 2. Testability
- Commands are simple data structures
- Handlers can be tested independently
- Easy to mock the bus for testing

### 3. Extensibility
- Add new commands without changing existing code
- Middleware can be added/removed easily
- Multiple handlers can be registered

### 4. Observability
- Centralized logging through middleware
- Easy to add metrics, tracing, etc.
- Command execution timing

## Adding a New Command

### Step 1: Define the Command

```go
// internal/mymodule/syncing/command.go
package syncing

import "akeneo-migrator/kit/bus"

const SyncMyEntityCommandType bus.Type = "myentity.sync"

type SyncMyEntityCommand struct {
    ID    string
    Debug bool
}

func (c SyncMyEntityCommand) Type() bus.Type {
    return SyncMyEntityCommandType
}
```

### Step 2: Create the Handler

```go
// internal/mymodule/syncing/command_handler.go
package syncing

import (
    "context"
    "akeneo-migrator/kit/bus"
)

type CommandHandler struct {
    service *Service
}

func NewCommandHandler(service *Service) *CommandHandler {
    return &CommandHandler{service: service}
}

func (h *CommandHandler) Handle(ctx context.Context, msg bus.Message) (bus.Response, error) {
    cmd := msg.(SyncMyEntityCommand)
    result, err := h.service.Sync(ctx, cmd.ID)
    return bus.Response{Data: result}, err
}
```

### Step 3: Register in Bootstrap

```go
// cmd/app/bootstrap/bootstrap.go

// Create service
myService := mymodule_syncing.NewService(sourceRepo, destRepo)

// Register handler
commandBus.Register(
    mymodule_syncing.SyncMyEntityCommandType,
    mymodule_syncing.NewCommandHandler(myService),
)
```

### Step 4: Dispatch from CLI

```go
response, err := app.CommandBus.Dispatch(ctx, mymodule_syncing.SyncMyEntityCommand{
    ID:    entityID,
    Debug: debug,
})
```

## Creating Custom Middleware

```go
// kit/bus/in_memory/middleware/metrics.go
package middleware

import (
    "context"
    "time"
    "akeneo-migrator/kit/bus"
    "akeneo-migrator/kit/bus/in_memory"
)

func Metrics() inmemory.Middleware {
    return func(ctx context.Context, msg bus.Message, next inmemory.NextFunc) (bus.Response, error) {
        start := time.Now()
        
        response, err := next(ctx, msg)
        
        duration := time.Since(start)
        
        // Record metrics
        recordCommandMetric(msg.Type(), duration, err)
        
        return response, err
    }
}
```

Then add it to the bus:

```go
commandBus := inmemory.NewCommandBus(
    middleware.Logging(),
    middleware.Metrics(),
)
```

## Testing

### Testing Commands

```go
func TestSyncAttributeCommand(t *testing.T) {
    cmd := SyncAttributeCommand{
        Code:  "sku",
        Debug: true,
    }
    
    if cmd.Type() != SyncAttributeCommandType {
        t.Error("Wrong command type")
    }
}
```

### Testing Handlers

```go
func TestCommandHandler(t *testing.T) {
    mockService := &MockService{}
    handler := NewCommandHandler(mockService)
    
    cmd := SyncAttributeCommand{Code: "sku"}
    response, err := handler.Handle(context.Background(), cmd)
    
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}
```

### Testing with Mock Bus

```go
type MockBus struct {
    dispatchFunc func(ctx context.Context, msg bus.Message) (bus.Response, error)
}

func (m *MockBus) Dispatch(ctx context.Context, msg bus.Message) (bus.Response, error) {
    return m.dispatchFunc(ctx, msg)
}
```

## Best Practices

1. **Keep Commands Simple**: Commands should only contain data, no logic
2. **One Handler per Command**: Each command type should have exactly one handler
3. **Immutable Commands**: Commands should be immutable once created
4. **Error Handling**: Always check response.Error in addition to returned error
5. **Type Safety**: Use type assertions carefully when extracting response data
6. **Middleware Order**: Order matters - logging should typically be first
7. **Context Propagation**: Always pass context through the chain

## Future Enhancements

Possible improvements:
- Async command execution
- Command queuing
- Retry middleware
- Circuit breaker middleware
- Command validation middleware
- Distributed tracing integration
- Event sourcing support
