package syncing

import (
	"context"
	"fmt"

	"akeneo-migrator/internal/category"
)

// Service handles category synchronization
type Service struct {
	sourceRepo category.SourceRepository
	destRepo   category.DestRepository
}

// NewService creates a new category sync service
func NewService(sourceRepo category.SourceRepository, destRepo category.DestRepository) *Service {
	return &Service{
		sourceRepo: sourceRepo,
		destRepo:   destRepo,
	}
}

// SyncResult contains the result of a sync operation
type SyncResult struct {
	Code    string
	Success bool
	Error   string
}

// Sync synchronizes a single category from source to destination
func (s *Service) Sync(ctx context.Context, code string) (*SyncResult, error) {
	result := &SyncResult{
		Code:    code,
		Success: false,
	}

	// 1. Get category from source
	categoryData, err := s.sourceRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error fetching category from source: %w", err)
	}

	// 2. Save category to destination
	err = s.destRepo.Save(ctx, code, categoryData)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, fmt.Errorf("error saving category to destination: %w", err)
	}

	result.Success = true
	return result, nil
}
