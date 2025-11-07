package syncing

import (
	"context"
	"fmt"

	"akeneo-migrator/internal/reference_entity"
)

// Service handles the synchronization logic for Reference Entities
type Service struct {
	sourceRepo reference_entity.SourceRepository
	destRepo   reference_entity.DestRepository
}

// NewService creates a new instance of the synchronization service
func NewService(sourceRepo reference_entity.SourceRepository, destRepo reference_entity.DestRepository) *Service {
	return &Service{
		sourceRepo: sourceRepo,
		destRepo:   destRepo,
	}
}

// SyncResult contains the result of a synchronization operation
type SyncResult struct {
	EntityName   string
	TotalRecords int
	SuccessCount int
	ErrorCount   int
	Errors       []SyncError
}

// SyncError represents an error during synchronization
type SyncError struct {
	Code    string
	Message string
}

// Sync synchronizes all records from a Reference Entity from source to destination
func (s *Service) Sync(ctx context.Context, entityName string) (*SyncResult, error) {
	result := &SyncResult{
		EntityName: entityName,
		Errors:     make([]SyncError, 0),
	}

	// 1. Get all records from source
	records, err := s.sourceRepo.FindAll(ctx, entityName)
	if err != nil {
		return nil, fmt.Errorf("error fetching records from source: %w", err)
	}

	result.TotalRecords = len(records)

	// 2. Sync each record to destination
	for _, record := range records {
		code, ok := record["code"].(string)
		if !ok {
			result.ErrorCount++
			result.Errors = append(result.Errors, SyncError{
				Code:    "unknown",
				Message: "could not extract record code",
			})
			continue
		}

		err := s.destRepo.Save(ctx, entityName, code, record)
		if err != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, SyncError{
				Code:    code,
				Message: err.Error(),
			})
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}
