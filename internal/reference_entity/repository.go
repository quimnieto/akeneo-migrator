package reference_entity

import "context"

// Record represents a Reference Entity record
type Record map[string]interface{}

// SourceRepository defines read-only operations for the source
type SourceRepository interface {
	// FindAll retrieves all records from a Reference Entity
	FindAll(ctx context.Context, entityName string) ([]Record, error)
}

// DestRepository defines read and write operations for the destination
type DestRepository interface {
	// FindAll retrieves all records from a Reference Entity
	FindAll(ctx context.Context, entityName string) ([]Record, error)
	
	// Save creates or updates a record in a Reference Entity
	Save(ctx context.Context, entityName string, code string, record Record) error
}
