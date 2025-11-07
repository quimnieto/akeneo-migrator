package syncing_test

import (
	"context"
	"errors"
	"testing"

	"akeneo-migrator/internal/reference_entity"
	"akeneo-migrator/internal/reference_entity/syncing"
)

// MockSourceRepository es un mock del repository de origen para testing
type MockSourceRepository struct {
	findAllFunc func(ctx context.Context, entityName string) ([]reference_entity.Record, error)
}

func (m *MockSourceRepository) FindAll(ctx context.Context, entityName string) ([]reference_entity.Record, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx, entityName)
	}
	return nil, nil
}

// MockDestRepository es un mock del repository de destino para testing
type MockDestRepository struct {
	findAllFunc func(ctx context.Context, entityName string) ([]reference_entity.Record, error)
	saveFunc    func(ctx context.Context, entityName string, code string, record reference_entity.Record) error
}

func (m *MockDestRepository) FindAll(ctx context.Context, entityName string) ([]reference_entity.Record, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx, entityName)
	}
	return nil, nil
}

func (m *MockDestRepository) Save(ctx context.Context, entityName string, code string, record reference_entity.Record) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, entityName, code, record)
	}
	return nil
}

func TestSync_Success(t *testing.T) {
	// Arrange
	mockRecords := []reference_entity.Record{
		{"code": "record1", "label": "Label 1"},
		{"code": "record2", "label": "Label 2"},
	}

	sourceRepo := &MockSourceRepository{
		findAllFunc: func(ctx context.Context, entityName string) ([]reference_entity.Record, error) {
			return mockRecords, nil
		},
	}

	destRepo := &MockDestRepository{
		saveFunc: func(ctx context.Context, entityName string, code string, record reference_entity.Record) error {
			return nil
		},
	}

	service := syncing.NewService(sourceRepo, destRepo)

	// Act
	result, err := service.Sync(context.Background(), "test_entity")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.TotalRecords != 2 {
		t.Errorf("Expected 2 total records, got %d", result.TotalRecords)
	}

	if result.SuccessCount != 2 {
		t.Errorf("Expected 2 successful records, got %d", result.SuccessCount)
	}

	if result.ErrorCount != 0 {
		t.Errorf("Expected 0 errors, got %d", result.ErrorCount)
	}
}

func TestSync_WithErrors(t *testing.T) {
	// Arrange
	mockRecords := []reference_entity.Record{
		{"code": "record1", "label": "Label 1"},
		{"code": "record2", "label": "Label 2"},
	}

	sourceRepo := &MockSourceRepository{
		findAllFunc: func(ctx context.Context, entityName string) ([]reference_entity.Record, error) {
			return mockRecords, nil
		},
	}

	destRepo := &MockDestRepository{
		saveFunc: func(ctx context.Context, entityName string, code string, record reference_entity.Record) error {
			if code == "record2" {
				return errors.New("error saving record2")
			}
			return nil
		},
	}

	service := syncing.NewService(sourceRepo, destRepo)

	// Act
	result, err := service.Sync(context.Background(), "test_entity")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.SuccessCount != 1 {
		t.Errorf("Expected 1 successful record, got %d", result.SuccessCount)
	}

	if result.ErrorCount != 1 {
		t.Errorf("Expected 1 error, got %d", result.ErrorCount)
	}

	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error in list, got %d", len(result.Errors))
	}
}

func TestSync_SourceError(t *testing.T) {
	// Arrange
	sourceRepo := &MockSourceRepository{
		findAllFunc: func(ctx context.Context, entityName string) ([]reference_entity.Record, error) {
			return nil, errors.New("source connection error")
		},
	}

	destRepo := &MockDestRepository{}

	service := syncing.NewService(sourceRepo, destRepo)

	// Act
	_, err := service.Sync(context.Background(), "test_entity")

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
