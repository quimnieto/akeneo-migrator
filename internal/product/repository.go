package product

import "context"

// Product represents a product
type Product map[string]interface{}

// ProductModel represents a product model
type ProductModel map[string]interface{}

// SourceRepository defines read-only operations for the source
type SourceRepository interface {
	// FindByIdentifier retrieves a product by its identifier
	FindByIdentifier(ctx context.Context, identifier string) (Product, error)
	
	// FindModelByCode retrieves a product model by its code
	FindModelByCode(ctx context.Context, code string) (ProductModel, error)
	
	// FindProductsByParent retrieves all products with a specific parent
	FindProductsByParent(ctx context.Context, parentCode string) ([]Product, error)
	
	// FindModelsByParent retrieves all product models with a specific parent
	FindModelsByParent(ctx context.Context, parentCode string) ([]ProductModel, error)
}

// DestRepository defines read and write operations for the destination
type DestRepository interface {
	// FindByIdentifier retrieves a product by its identifier
	FindByIdentifier(ctx context.Context, identifier string) (Product, error)
	
	// Save creates or updates a product
	Save(ctx context.Context, identifier string, product Product) error
	
	// FindModelByCode retrieves a product model by its code
	FindModelByCode(ctx context.Context, code string) (ProductModel, error)
	
	// SaveModel creates or updates a product model
	SaveModel(ctx context.Context, code string, model ProductModel) error
	
	// FindProductsByParent retrieves all products with a specific parent
	FindProductsByParent(ctx context.Context, parentCode string) ([]Product, error)
	
	// FindModelsByParent retrieves all product models with a specific parent
	FindModelsByParent(ctx context.Context, parentCode string) ([]ProductModel, error)
}
