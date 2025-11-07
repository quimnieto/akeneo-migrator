package reference_entity

import "context"

// Record represents a Reference Entity record
type Record map[string]interface{}

// Entity represents a Reference Entity definition
type Entity map[string]interface{}

// SourceRepository defines read-only operations for the source
type SourceRepository interface {
	// FindEntity retrieves a Reference Entity definition
	FindEntity(ctx context.Context, entityCode string) (Entity, error)
	
	// FindAll retrieves all records from a Reference Entity
	FindAll(ctx context.Context, entityName string) ([]Record, error)
}

// DestRepository defines read and write operations for the destination
type DestRepository interface {
	// FindEntity retrieves a Reference Entity definition
	FindEntity(ctx context.Context, entityCode string) (Entity, error)
	
	// SaveEntity creates or updates a Reference Entity definition
	SaveEntity(ctx context.Context, entityCode string, entity Entity) error
	
	// FindAll retrieves all records from a Reference Entity
	FindAll(ctx context.Context, entityName string) ([]Record, error)
	
	// Save creates or updates a record in a Reference Entity
	Save(ctx context.Context, entityName string, code string, record Record) error
}
