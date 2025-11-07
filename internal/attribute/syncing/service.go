package syncing

import (
	"context"
	"fmt"

	"akeneo-migrator/internal/attribute"
)

// Service handles attribute synchronization
type Service struct {
	sourceRepo attribute.SourceRepository
	destRepo   attribute.DestRepository
}

// NewService creates a new attribute sync service
func NewService(sourceRepo attribute.SourceRepository, destRepo attribute.DestRepository) *Service {
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

// Sync synchronizes a single attribute from source to destination
func (s *Service) Sync(ctx context.Context, code string) (*SyncResult, error) {
	result := &SyncResult{
		Code:    code,
		Success: false,
	}

	// 1. Get attribute from source
	attributeData, err := s.sourceRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error fetching attribute from source: %w", err)
	}

	// 2. Save attribute to destination
	err = s.destRepo.Save(ctx, code, attributeData)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, fmt.Errorf("error saving attribute to destination: %w", err)
	}

	result.Success = true
	return result, nil
}
