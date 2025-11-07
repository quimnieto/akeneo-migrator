package attribute

import "context"

// Attribute represents an attribute
type Attribute map[string]interface{}

// SourceRepository defines read-only operations for attributes from source
type SourceRepository interface {
	// FindByCode retrieves an attribute by its code
	FindByCode(ctx context.Context, code string) (Attribute, error)
}

// DestRepository defines write operations for attributes to destination
type DestRepository interface {
	// Save creates or updates an attribute
	Save(ctx context.Context, code string, attribute Attribute) error
}
