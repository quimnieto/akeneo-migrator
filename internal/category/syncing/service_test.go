package syncing

import (
	"context"
	"errors"
	"testing"

	"akeneo-migrator/internal/category"
)

// Mock repositories
type mockSourceRepo struct {
	findByCodeFunc func(ctx context.Context, code string) (category.Category, error)
}

func (m *mockSourceRepo) FindByCode(ctx context.Context, code string) (category.Category, error) {
	if m.findByCodeFunc != nil {
		return m.findByCodeFunc(ctx, code)
	}
	return nil, nil
}

type mockDestRepo struct {
	saveFunc func(ctx context.Context, code string, cat category.Category) error
}

func (m *mockDestRepo) Save(ctx context.Context, code string, cat category.Category) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, code, cat)
	}
	return nil
}

func TestSync_Success(t *testing.T) {
	sourceRepo := &mockSourceRepo{
		findByCodeFunc: func(ctx context.Context, code string) (category.Category, error) {
			return category.Category{
				"code":   code,
				"parent": "master",
				"labels": map[string]string{"en_US": "Test Category"},
			}, nil
		},
	}

	destRepo := &mockDestRepo{
		saveFunc: func(ctx context.Context, code string, cat category.Category) error {
			return nil
		},
	}

	service := NewService(sourceRepo, destRepo)
	result, err := service.Sync(context.Background(), "test_category")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result.Success {
		t.Error("Expected success to be true")
	}

	if result.Code != "test_category" {
		t.Errorf("Expected code 'test_category', got '%s'", result.Code)
	}
}

func TestSync_SourceError(t *testing.T) {
	sourceRepo := &mockSourceRepo{
		findByCodeFunc: func(ctx context.Context, code string) (category.Category, error) {
			return nil, errors.New("source error")
		},
	}

	destRepo := &mockDestRepo{}

	service := NewService(sourceRepo, destRepo)
	_, err := service.Sync(context.Background(), "test_category")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestSync_DestError(t *testing.T) {
	sourceRepo := &mockSourceRepo{
		findByCodeFunc: func(ctx context.Context, code string) (category.Category, error) {
			return category.Category{
				"code":   code,
				"parent": "master",
			}, nil
		},
	}

	destRepo := &mockDestRepo{
		saveFunc: func(ctx context.Context, code string, cat category.Category) error {
			return errors.New("dest error")
		},
	}

	service := NewService(sourceRepo, destRepo)
	result, err := service.Sync(context.Background(), "test_category")

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result.Success {
		t.Error("Expected success to be false")
	}

	if result.Error == "" {
		t.Error("Expected error message in result")
	}
}
