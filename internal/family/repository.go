package family

import "context"

// Family represents a family
type Family map[string]interface{}

// FamilyVariant represents a family variant
type FamilyVariant map[string]interface{}

// SourceRepository defines read-only operations for families from source
type SourceRepository interface {
	// FindByCode retrieves a family by its code
	FindByCode(ctx context.Context, code string) (Family, error)

	// GetVariants retrieves all variants for a family
	GetVariants(ctx context.Context, familyCode string) ([]FamilyVariant, error)
}

// DestRepository defines write operations for families to destination
type DestRepository interface {
	// Save creates or updates a family
	Save(ctx context.Context, code string, family Family) error

	// SaveVariant creates or updates a family variant
	SaveVariant(ctx context.Context, familyCode, variantCode string, variant FamilyVariant) error
}
