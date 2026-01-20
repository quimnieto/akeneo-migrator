package attribute

import "context"

// Attribute represents an attribute
type Attribute map[string]interface{}

// AttributeOption represents an attribute option
type AttributeOption map[string]interface{}

// SourceRepository defines read-only operations for attributes from source
type SourceRepository interface {
	// FindByCode retrieves an attribute by its code
	FindByCode(ctx context.Context, code string) (Attribute, error)

	// GetOptions retrieves all options for an attribute
	GetOptions(ctx context.Context, attributeCode string) ([]AttributeOption, error)
}

// DestRepository defines write operations for attributes to destination
type DestRepository interface {
	// Save creates or updates an attribute
	Save(ctx context.Context, code string, attribute Attribute) error

	// SaveOption creates or updates an attribute option
	SaveOption(ctx context.Context, attributeCode, optionCode string, option AttributeOption) error
}
