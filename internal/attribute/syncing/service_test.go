package syncing

import (
	"context"
	"errors"
	"testing"

	"akeneo-migrator/internal/attribute"
)

// Mock repositories
type mockSourceRepo struct {
	findByCodeFunc func(ctx context.Context, code string) (attribute.Attribute, error)
}

func (m *mockSourceRepo) FindByCode(ctx context.Context, code string) (attribute.Attribute, error) {
	if m.findByCodeFunc != nil {
		return m.findByCodeFunc(ctx, code)
	}
	return nil, nil
}

type mockDestRepo struct {
	saveFunc func(ctx context.Context, code string, attr attribute.Attribute) error
}

func (m *mockDestRepo) Save(ctx context.Context, code string, attr attribute.Attribute) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, code, attr)
	}
	return nil
}

func TestSync_Success(t *testing.T) {
	sourceRepo := &mockSourceRepo{
		findByCodeFunc: func(ctx context.Context, code string) (attribute.Attribute, error) {
			return attribute.Attribute{
				"code": code,
				"type": "pim_catalog_text",
			}, nil
		},
	}

	destRepo := &mockDestRepo{
		saveFunc: func(ctx context.Context, code string, attr attribute.Attribute) error {
			return nil
		},
	}

	service := NewService(sourceRepo, destRepo)
	result, err := service.Sync(context.Background(), "sku")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result.Success {
		t.Error("Expected success to be true")
	}

	if result.Code != "sku" {
		t.Errorf("Expected code 'sku', got '%s'", result.Code)
	}
}

func TestSync_SourceError(t *testing.T) {
	sourceRepo := &mockSourceRepo{
		findByCodeFunc: func(ctx context.Context, code string) (attribute.Attribute, error) {
			return nil, errors.New("source error")
		},
	}

	destRepo := &mockDestRepo{}

	service := NewService(sourceRepo, destRepo)
	_, err := service.Sync(context.Background(), "sku")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestSync_DestError(t *testing.T) {
	sourceRepo := &mockSourceRepo{
		findByCodeFunc: func(ctx context.Context, code string) (attribute.Attribute, error) {
			return attribute.Attribute{
				"code": code,
				"type": "pim_catalog_text",
			}, nil
		},
	}

	destRepo := &mockDestRepo{
		saveFunc: func(ctx context.Context, code string, attr attribute.Attribute) error {
			return errors.New("dest error")
		},
	}

	service := NewService(sourceRepo, destRepo)
	result, err := service.Sync(context.Background(), "sku")

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
