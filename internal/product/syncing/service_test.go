package syncing_test

import (
	"context"
	"errors"
	"testing"

	"akeneo-migrator/internal/product"
	"akeneo-migrator/internal/product/syncing"
)

// MockSourceRepository is a mock of the source repository for testing
type MockSourceRepository struct {
	findByIdentifierFunc     func(ctx context.Context, identifier string) (product.Product, error)
	findModelByCodeFunc      func(ctx context.Context, code string) (product.ProductModel, error)
	findProductsByParentFunc func(ctx context.Context, parentCode string) ([]product.Product, error)
	findModelsByParentFunc   func(ctx context.Context, parentCode string) ([]product.ProductModel, error)
}

func (m *MockSourceRepository) FindByIdentifier(ctx context.Context, identifier string) (product.Product, error) {
	if m.findByIdentifierFunc != nil {
		return m.findByIdentifierFunc(ctx, identifier)
	}
	return product.Product{"identifier": identifier}, nil
}

func (m *MockSourceRepository) FindModelByCode(ctx context.Context, code string) (product.ProductModel, error) {
	if m.findModelByCodeFunc != nil {
		return m.findModelByCodeFunc(ctx, code)
	}
	return product.ProductModel{"code": code}, nil
}

func (m *MockSourceRepository) FindProductsByParent(ctx context.Context, parentCode string) ([]product.Product, error) {
	if m.findProductsByParentFunc != nil {
		return m.findProductsByParentFunc(ctx, parentCode)
	}
	return []product.Product{}, nil
}

func (m *MockSourceRepository) FindModelsByParent(ctx context.Context, parentCode string) ([]product.ProductModel, error) {
	if m.findModelsByParentFunc != nil {
		return m.findModelsByParentFunc(ctx, parentCode)
	}
	return []product.ProductModel{}, nil
}

// MockDestRepository is a mock of the destination repository for testing
type MockDestRepository struct {
	findByIdentifierFunc     func(ctx context.Context, identifier string) (product.Product, error)
	saveFunc                 func(ctx context.Context, identifier string, productData product.Product) error
	findModelByCodeFunc      func(ctx context.Context, code string) (product.ProductModel, error)
	saveModelFunc            func(ctx context.Context, code string, model product.ProductModel) error
	findProductsByParentFunc func(ctx context.Context, parentCode string) ([]product.Product, error)
	findModelsByParentFunc   func(ctx context.Context, parentCode string) ([]product.ProductModel, error)
}

func (m *MockDestRepository) FindByIdentifier(ctx context.Context, identifier string) (product.Product, error) {
	if m.findByIdentifierFunc != nil {
		return m.findByIdentifierFunc(ctx, identifier)
	}
	return product.Product{"identifier": identifier}, nil
}

func (m *MockDestRepository) Save(ctx context.Context, identifier string, productData product.Product) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, identifier, productData)
	}
	return nil
}

func (m *MockDestRepository) FindModelByCode(ctx context.Context, code string) (product.ProductModel, error) {
	if m.findModelByCodeFunc != nil {
		return m.findModelByCodeFunc(ctx, code)
	}
	return product.ProductModel{"code": code}, nil
}

func (m *MockDestRepository) SaveModel(ctx context.Context, code string, model product.ProductModel) error {
	if m.saveModelFunc != nil {
		return m.saveModelFunc(ctx, code, model)
	}
	return nil
}

func (m *MockDestRepository) FindProductsByParent(ctx context.Context, parentCode string) ([]product.Product, error) {
	if m.findProductsByParentFunc != nil {
		return m.findProductsByParentFunc(ctx, parentCode)
	}
	return []product.Product{}, nil
}

func (m *MockDestRepository) FindModelsByParent(ctx context.Context, parentCode string) ([]product.ProductModel, error) {
	if m.findModelsByParentFunc != nil {
		return m.findModelsByParentFunc(ctx, parentCode)
	}
	return []product.ProductModel{}, nil
}

func TestSync_Success(t *testing.T) {
	// Arrange
	mockProduct := product.Product{
		"identifier": "SKU-123",
		"family":     "shoes",
		"enabled":    true,
	}

	sourceRepo := &MockSourceRepository{
		findByIdentifierFunc: func(ctx context.Context, identifier string) (product.Product, error) {
			return mockProduct, nil
		},
	}

	destRepo := &MockDestRepository{
		saveFunc: func(ctx context.Context, identifier string, productData product.Product) error {
			return nil
		},
	}

	service := syncing.NewService(sourceRepo, destRepo)

	// Act
	result, err := service.Sync(context.Background(), "SKU-123")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success to be true, got false")
	}

	if result.Identifier != "SKU-123" {
		t.Errorf("Expected identifier 'SKU-123', got %s", result.Identifier)
	}
}

func TestSync_SourceError(t *testing.T) {
	// Arrange
	sourceRepo := &MockSourceRepository{
		findByIdentifierFunc: func(ctx context.Context, identifier string) (product.Product, error) {
			return nil, errors.New("product not found")
		},
	}

	destRepo := &MockDestRepository{}

	service := syncing.NewService(sourceRepo, destRepo)

	// Act
	_, err := service.Sync(context.Background(), "SKU-123")

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestSync_DestError(t *testing.T) {
	// Arrange
	mockProduct := product.Product{
		"identifier": "SKU-123",
		"family":     "shoes",
	}

	sourceRepo := &MockSourceRepository{
		findByIdentifierFunc: func(ctx context.Context, identifier string) (product.Product, error) {
			return mockProduct, nil
		},
	}

	destRepo := &MockDestRepository{
		saveFunc: func(ctx context.Context, identifier string, productData product.Product) error {
			return errors.New("save failed")
		},
	}

	service := syncing.NewService(sourceRepo, destRepo)

	// Act
	result, err := service.Sync(context.Background(), "SKU-123")

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result.Success {
		t.Error("Expected success to be false, got true")
	}
}
