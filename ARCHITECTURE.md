# Architecture Documentation

## Overview

This project implements a **Hexagonal Architecture** (Ports & Adapters) with a **Command Bus** pattern for executing synchronization operations between Akeneo PIM instances.

## Directory Structure

```
akeneo-migrator/
├── cmd/
│   └── app/
│       ├── main.go                    # Application entry point
│       └── bootstrap/
│           └── bootstrap.go           # Dependency injection & CLI setup
│
├── internal/
│   ├── attribute/                     # Attribute domain
│   │   ├── repository.go              # Domain interfaces
│   │   └── syncing/
│   │       ├── service.go             # Business logic
│   │       ├── command.go             # Command definition
│   │       ├── command_handler.go     # Command handler
│   │       └── service_test.go        # Unit tests
│   │
│   ├── product/                       # Product domain
│   │   ├── repository.go
│   │   └── syncing/
│   │       ├── service.go
│   │       ├── command.go
│   │       ├── command_handler.go
│   │       └── service_test.go
│   │
│   ├── reference_entity/              # Reference Entity domain
│   │   ├── repository.go
│   │   └── syncing/
│   │       ├── service.go
│   │       ├── command.go
│   │       ├── command_handler.go
│   │       └── service_test.go
│   │
│   └── platform/                      # Infrastructure layer
│       ├── client/
│       │   └── akeneo/
│       │       └── client.go          # HTTP client for Akeneo API
│       ├── storage/
│       │   └── akeneo/
│       │       ├── attribute_repository.go
│       │       ├── product_repository.go
│       │       └── reference_entity_repository.go
│       └── config/
│           └── config.go              # Configuration management
│
├── kit/                               # Shared kernel / reusable components
│   ├── bus/
│   │   ├── bus.go                     # Bus interfaces
│   │   └── in_memory/
│   │       ├── command_bus.go         # In-memory implementation
│   │       ├── command_bus_test.go
│   │       └── middleware/
│   │           └── logging.go         # Logging middleware
│   └── config/
│       └── static/
│           ├── config.go              # Config interfaces
│           └── viper/
│               └── viper.go           # Viper implementation
│
├── configs/                           # Configuration files
│   └── akeneo-migrator/
│       └── settings.local.json
│
├── Makefile                           # Build automation
├── go.mod                             # Go dependencies
├── .golangci.yml                      # Linter configuration
└── .github/
    └── workflows/
        └── ci.yml                     # CI/CD pipeline
```

## Architectural Layers

### 1. Domain Layer (`internal/[module]/`)

Contains the core business logic and domain interfaces (ports).

**Responsibilities:**
- Define domain entities and value objects
- Define repository interfaces (ports)
- Pure business logic, no infrastructure dependencies

**Example:**
```go
// internal/attribute/repository.go
package attribute

type Attribute map[string]interface{}

type SourceRepository interface {
    FindByCode(ctx context.Context, code string) (Attribute, error)
}

type DestRepository interface {
    Save(ctx context.Context, code string, attribute Attribute) error
}
```

### 2. Application Layer (`internal/[module]/syncing/`)

Contains use cases and application-specific business rules.

**Components:**
- **Services**: Orchestrate domain logic
- **Commands**: Define intentions (what to do)
- **Handlers**: Execute commands (how to do it)

**Example:**
```go
// Service
type Service struct {
    sourceRepo SourceRepository
    destRepo   DestRepository
}

// Command
type SyncAttributeCommand struct {
    Code  string
    Debug bool
}

// Handler
type CommandHandler struct {
    service *Service
}
```

### 3. Infrastructure Layer (`internal/platform/`)

Contains concrete implementations of domain interfaces (adapters).

**Components:**
- **Client**: HTTP communication with external APIs
- **Storage**: Repository implementations
- **Config**: Configuration management

**Example:**
```go
// internal/platform/storage/akeneo/attribute_repository.go
type SourceAttributeRepository struct {
    client *akeneo.Client
}

func (r *SourceAttributeRepository) FindByCode(ctx context.Context, code string) (attribute.Attribute, error) {
    return r.client.GetAttribute(code)
}
```

### 4. Shared Kernel (`kit/`)

Reusable components shared across the application.

**Components:**
- **Bus**: Command bus implementation
- **Config**: Configuration loaders
- **Middleware**: Cross-cutting concerns

### 5. Bootstrap (`cmd/app/bootstrap/`)

Wires everything together through dependency injection.

**Responsibilities:**
- Create instances of all components
- Register command handlers in the bus
- Set up CLI commands
- Configure middleware

## Data Flow

### Command Execution Flow

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI Layer                            │
│  (User runs: ./akeneo-migrator sync-attribute sku)          │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Bootstrap Layer                           │
│  • Parse CLI arguments                                       │
│  • Create SyncAttributeCommand                               │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                     Command Bus                              │
│  • Dispatch(ctx, SyncAttributeCommand)                       │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   Middleware Chain                           │
│  • Logging (before)                                          │
│  • [Future: Metrics, Tracing, Retry, etc.]                  │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   Command Handler                            │
│  • Extract command data                                      │
│  • Call service method                                       │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      Service                                 │
│  • Business logic                                            │
│  • Orchestrate repositories                                  │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   Repositories                               │
│  • SourceRepository.FindByCode()                             │
│  • DestRepository.Save()                                     │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   HTTP Client                                │
│  • GET /api/rest/v1/attributes/{code}                        │
│  • PATCH /api/rest/v1/attributes/{code}                      │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   Akeneo API                                 │
└─────────────────────────────────────────────────────────────┘
```

## Design Patterns

### 1. Hexagonal Architecture (Ports & Adapters)

**Benefits:**
- Business logic independent of external systems
- Easy to test with mocks
- Can swap implementations (e.g., different PIM systems)

**Implementation:**
- **Ports**: Interfaces in domain layer (`repository.go`)
- **Adapters**: Implementations in infrastructure layer (`storage/`, `client/`)

### 2. Command Bus Pattern

**Benefits:**
- Decouples command creation from execution
- Centralized middleware for cross-cutting concerns
- Easy to add new commands without changing existing code

**Implementation:**
- Commands: Simple data structures
- Handlers: Execute business logic
- Bus: Routes commands to handlers
- Middleware: Wraps execution

### 3. Repository Pattern

**Benefits:**
- Abstracts data access
- Separates read (Source) from write (Dest) operations
- Easy to mock for testing

**Implementation:**
- `SourceRepository`: Read-only operations
- `DestRepository`: Write operations

### 4. Dependency Injection

**Benefits:**
- Loose coupling
- Easy to test
- Flexible configuration

**Implementation:**
- All dependencies injected in `bootstrap.go`
- No global state or singletons

## Testing Strategy

### Unit Tests

Test individual components in isolation:

```go
// Service tests with mock repositories
func TestSync_Success(t *testing.T) {
    mockSource := &MockSourceRepo{}
    mockDest := &MockDestRepo{}
    service := NewService(mockSource, mockDest)
    // ...
}
```

### Integration Tests

Test component interactions:

```go
// Command bus with real handlers
func TestCommandBus_Integration(t *testing.T) {
    bus := inmemory.NewCommandBus()
    handler := NewCommandHandler(service)
    bus.Register(SyncCommandType, handler)
    // ...
}
```

### End-to-End Tests

Test complete flows (manual or automated):

```bash
# Test actual sync with test instances
./akeneo-migrator sync-attribute test_attribute
```

## Configuration Management

### Configuration Layers

1. **JSON Files** (`configs/akeneo-migrator/settings.local.json`)
2. **Environment Variables** (override JSON)
3. **Command-line Flags** (override environment)

### Configuration Flow

```
JSON File → Viper → Config Struct → Services
```

## Error Handling

### Error Propagation

Errors flow up through layers:

```
API Error → Client → Repository → Service → Handler → Bus → CLI
```

### Error Types

1. **Domain Errors**: Business rule violations
2. **Infrastructure Errors**: Network, API, database errors
3. **Validation Errors**: Invalid input data

## Security Considerations

1. **Credentials**: Stored in config files (not in code)
2. **OAuth2**: Token-based authentication with refresh
3. **HTTPS**: All API communication encrypted
4. **Secrets**: Should use environment variables in production

## Performance Considerations

1. **Pagination**: Large datasets fetched in chunks
2. **Connection Pooling**: HTTP client reuses connections
3. **Token Caching**: OAuth tokens cached until expiry
4. **Batch Operations**: Multiple items synced in single session

## Future Enhancements

### Planned Features

1. **Async Processing**: Queue-based command execution
2. **Retry Logic**: Automatic retry with exponential backoff
3. **Metrics**: Prometheus metrics for monitoring
4. **Distributed Tracing**: OpenTelemetry integration
5. **Event Sourcing**: Audit log of all operations
6. **Webhooks**: Real-time sync triggers
7. **Conflict Resolution**: Handle concurrent modifications

### Scalability

1. **Horizontal Scaling**: Multiple workers processing commands
2. **Message Queue**: RabbitMQ or Kafka for command distribution
3. **Caching**: Redis for frequently accessed data
4. **Rate Limiting**: Respect API limits

## Contributing

When adding new features:

1. Follow the existing directory structure
2. Keep domain logic pure (no infrastructure dependencies)
3. Use command bus for all operations
4. Write tests for all layers
5. Update documentation

## References

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Command Bus Pattern](https://matthiasnoback.nl/2015/01/a-wave-of-command-buses/)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)
- [Dependency Injection](https://martinfowler.com/articles/injection.html)
