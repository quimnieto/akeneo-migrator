# Project Architecture

## Pattern: Hexagonal Architecture (Ports & Adapters)

This project implements a hexagonal architecture that clearly separates business logic from implementation details.

## Layers

### 1. Domain (`internal/reference_entity/`)

Contains pure business logic and interfaces (ports).

```go
// repository.go - Defines separate contracts

// SourceRepository - Read-only for source
type SourceRepository interface {
    FindAll(ctx context.Context, entityName string) ([]Record, error)
}

// DestRepository - Read and write for destination
type DestRepository interface {
    FindAll(ctx context.Context, entityName string) ([]Record, error)
    Save(ctx context.Context, entityName string, code string, record Record) error
}
```

**Features:**
- ✅ No external dependencies
- ✅ Easy to test
- ✅ Independent of infrastructure
- ✅ **ISP (Interface Segregation Principle)**: Separate interfaces by responsibility

### 2. Services (`internal/reference_entity/syncing/`)

Implements application use cases.

```go
// service.go - Synchronization logic
type Service struct {
    sourceRepo SourceRepository  // Read-only
    destRepo   DestRepository    // Read and write
}

func (s *Service) Sync(ctx context.Context, entityName string) (*SyncResult, error)
```

**Features:**
- ✅ Orchestrates domain operations
- ✅ Doesn't know implementation details
- ✅ Uses domain interfaces
- ✅ **Clear separation**: Source only reads, Dest reads and writes

### 3. Infrastructure (`internal/platform/`)

Implements concrete adapters.

#### HTTP Client (`internal/platform/client/akeneo/`)
```go
// client.go - HTTP client for Akeneo API
type Client struct {
    config      ClientConfig
    httpClient  *http.Client
    accessToken string
}
```

#### Repository (`internal/platform/storage/akeneo/`)
```go
// reference_entity_repository.go - Separate implementations

// SourceReferenceEntityRepository - Read-only
type SourceReferenceEntityRepository struct {
    client *akeneo.Client
}

func (r *SourceReferenceEntityRepository) FindAll(ctx context.Context, entityName string) ([]Record, error)

// DestReferenceEntityRepository - Read and write
type DestReferenceEntityRepository struct {
    client *akeneo.Client
}

func (r *DestReferenceEntityRepository) FindAll(ctx context.Context, entityName string) ([]Record, error)
func (r *DestReferenceEntityRepository) Save(ctx context.Context, entityName string, code string, record Record) error
```

**Features:**
- ✅ Implements domain interfaces
- ✅ Handles technical details (HTTP, JSON, authentication)
- ✅ Can be replaced without affecting business logic
- ✅ **Separation of concerns**: Source cannot modify data

### 4. Bootstrap (`cmd/app/bootstrap/`)

Dependency injection and configuration.

```go
// bootstrap.go
func Run() error {
    // 1. Load configuration
    // 2. Create clients
    // 3. Create repositories
    // 4. Create services
    // 5. Configure CLI commands
    // 6. Execute
}
```

**Features:**
- ✅ Single point of configuration
- ✅ Manual dependency injection
- ✅ Easy to understand and maintain

## Data Flow

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI Command                          │
│                    (cmd/app/bootstrap)                       │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Syncing Service                           │
│            (internal/reference_entity/syncing)               │
│                                                              │
│  - Orchestrates synchronization                             │
│  - Handles errors and results                               │
└──────────────┬──────────────────────────┬───────────────────┘
               │                          │
               ▼                          ▼
┌──────────────────────────┐  ┌──────────────────────────┐
│   Source Repository      │  │    Dest Repository       │
│  (interface)             │  │   (interface)            │
└──────────┬───────────────┘  └──────────┬───────────────┘
           │                              │
           ▼                              ▼
┌──────────────────────────┐  ┌──────────────────────────┐
│ Akeneo Repository Impl   │  │ Akeneo Repository Impl   │
│ (platform/storage)       │  │ (platform/storage)       │
└──────────┬───────────────┘  └──────────┬───────────────┘
           │                              │
           ▼                              ▼
┌──────────────────────────┐  ┌──────────────────────────┐
│   Akeneo HTTP Client     │  │   Akeneo HTTP Client     │
│   (platform/client)      │  │   (platform/client)      │
└──────────┬───────────────┘  └──────────┬───────────────┘
           │                              │
           ▼                              ▼
     Source Akeneo API            Dest Akeneo API
```

## Advantages of this Architecture

### 1. Testability
```go
// Easy to create mocks for testing
type MockSourceRepository struct {
    findAllFunc func(ctx context.Context, entityName string) ([]Record, error)
}
```

### 2. Decoupling
- Business logic doesn't depend on Akeneo
- You can change PIM without touching the domain
- Services only know interfaces

### 3. Extensibility
```go
// Adding a new adapter is simple
type OtherPIMRepository struct {
    client *otherpim.Client
}

func (r *OtherPIMRepository) FindAll(...) ([]Record, error) {
    // Implementation for another PIM
}
```

### 4. Maintainability
- Clear separation of responsibilities
- Each layer has a specific purpose
- Easy to understand and modify

## Testing

### Unit Tests
```bash
# Test services with mocks
go test ./internal/reference_entity/syncing/...
```

### Integration Tests
```bash
# Test with real Akeneo (requires configuration)
go test ./internal/platform/storage/akeneo/... -tags=integration
```

## Adding New Features

### 1. Define in Domain
```go
// internal/reference_entity/repository.go
type DestRepository interface {
    FindAll(ctx context.Context, entityName string) ([]Record, error)
    Save(ctx context.Context, entityName string, code string, record Record) error
    Delete(ctx context.Context, entityName string, code string) error // New
}
```

### 2. Create Service
```go
// internal/reference_entity/deleting/service.go
type Service struct {
    repo DestRepository
}

func (s *Service) Delete(ctx context.Context, entityName, code string) error {
    return s.repo.Delete(ctx, entityName, code)
}
```

### 3. Implement in Infrastructure
```go
// internal/platform/storage/akeneo/reference_entity_repository.go
func (r *DestReferenceEntityRepository) Delete(ctx context.Context, entityName string, code string) error {
    return r.client.DeleteReferenceEntityRecord(entityName, code)
}
```

### 4. Register in Bootstrap
```go
// cmd/app/bootstrap/bootstrap.go
deleteService := deleting.NewService(destRepository)
deleteCmd := createDeleteCommand(deleteService)
rootCmd.AddCommand(deleteCmd)
```

## Applied Principles

- **SOLID**: Each component has a single responsibility
- **DIP**: We depend on abstractions, not implementations
- **ISP**: Small and specific interfaces
- **OCP**: Open for extension, closed for modification
