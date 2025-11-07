package akeneo

import (
	"context"

	"akeneo-migrator/internal/platform/client/akeneo"
	"akeneo-migrator/internal/product"
)

// SourceProductRepository implements the read-only repository for the source
type SourceProductRepository struct {
	client *akeneo.Client
}

// NewSourceProductRepository creates a new instance of the source repository
func NewSourceProductRepository(client *akeneo.Client) *SourceProductRepository {
	return &SourceProductRepository{
		client: client,
	}
}

// FindByIdentifier retrieves a product by its identifier
func (r *SourceProductRepository) FindByIdentifier(ctx context.Context, identifier string) (product.Product, error) {
	productData, err := r.client.GetProduct(identifier)
	if err != nil {
		return nil, err
	}

	return product.Product(productData), nil
}

// FindModelByCode retrieves a product model by its code
func (r *SourceProductRepository) FindModelByCode(ctx context.Context, code string) (product.ProductModel, error) {
	model, err := r.client.GetProductModel(code)
	if err != nil {
		return nil, err
	}

	return product.ProductModel(model), nil
}

// FindProductsByParent retrieves all products with a specific parent
func (r *SourceProductRepository) FindProductsByParent(ctx context.Context, parentCode string) ([]product.Product, error) {
	products, err := r.client.GetProductsByParent(parentCode)
	if err != nil {
		return nil, err
	}

	result := make([]product.Product, len(products))
	for i, p := range products {
		result[i] = product.Product(p)
	}

	return result, nil
}

// FindModelsByParent retrieves all product models with a specific parent
func (r *SourceProductRepository) FindModelsByParent(ctx context.Context, parentCode string) ([]product.ProductModel, error) {
	models, err := r.client.GetProductModelsByParent(parentCode)
	if err != nil {
		return nil, err
	}

	result := make([]product.ProductModel, len(models))
	for i, m := range models {
		result[i] = product.ProductModel(m)
	}

	return result, nil
}

// DestProductRepository implements the read/write repository for the destination
type DestProductRepository struct {
	client *akeneo.Client
}

// NewDestProductRepository creates a new instance of the destination repository
func NewDestProductRepository(client *akeneo.Client) *DestProductRepository {
	return &DestProductRepository{
		client: client,
	}
}

// FindByIdentifier retrieves a product by its identifier
func (r *DestProductRepository) FindByIdentifier(ctx context.Context, identifier string) (product.Product, error) {
	productData, err := r.client.GetProduct(identifier)
	if err != nil {
		return nil, err
	}

	return product.Product(productData), nil
}

// Save creates or updates a product
func (r *DestProductRepository) Save(ctx context.Context, identifier string, productData product.Product) error {
	// Convert from product.Product to akeneo.Product
	akeneoProduct := akeneo.Product(productData)
	return r.client.PatchProduct(identifier, akeneoProduct)
}

// FindModelByCode retrieves a product model by its code
func (r *DestProductRepository) FindModelByCode(ctx context.Context, code string) (product.ProductModel, error) {
	model, err := r.client.GetProductModel(code)
	if err != nil {
		return nil, err
	}

	return product.ProductModel(model), nil
}

// SaveModel creates or updates a product model
func (r *DestProductRepository) SaveModel(ctx context.Context, code string, model product.ProductModel) error {
	// Convert from product.ProductModel to akeneo.ProductModel
	akeneoModel := akeneo.ProductModel(model)
	return r.client.PatchProductModel(code, akeneoModel)
}

// FindProductsByParent retrieves all products with a specific parent
func (r *DestProductRepository) FindProductsByParent(ctx context.Context, parentCode string) ([]product.Product, error) {
	products, err := r.client.GetProductsByParent(parentCode)
	if err != nil {
		return nil, err
	}

	result := make([]product.Product, len(products))
	for i, p := range products {
		result[i] = product.Product(p)
	}

	return result, nil
}

// FindModelsByParent retrieves all product models with a specific parent
func (r *DestProductRepository) FindModelsByParent(ctx context.Context, parentCode string) ([]product.ProductModel, error) {
	models, err := r.client.GetProductModelsByParent(parentCode)
	if err != nil {
		return nil, err
	}

	result := make([]product.ProductModel, len(models))
	for i, m := range models {
		result[i] = product.ProductModel(m)
	}

	return result, nil
}
