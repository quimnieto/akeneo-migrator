package category

import "context"

// Category represents a category
type Category map[string]interface{}

// SourceRepository defines read-only operations for categories from source
type SourceRepository interface {
	// FindByCode retrieves a category by its code
	FindByCode(ctx context.Context, code string) (Category, error)
}

// DestRepository defines write operations for categories to destination
type DestRepository interface {
	// Save creates or updates a category
	Save(ctx context.Context, code string, category Category) error
}
