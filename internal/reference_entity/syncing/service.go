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

// Sync synchronizes a Reference Entity (definition + attributes + records) from source to destination
func (s *Service) Sync(ctx context.Context, entityName string) (*SyncResult, error) {
	result := &SyncResult{
		EntityName: entityName,
		Errors:     make([]SyncError, 0),
	}

	// 1. Get Reference Entity definition from source
	entity, err := s.sourceRepo.FindEntity(ctx, entityName)
	if err != nil {
		return nil, fmt.Errorf("error fetching reference entity definition from source: %w", err)
	}

	// 2. Create or update Reference Entity in destination
	err = s.destRepo.SaveEntity(ctx, entityName, entity)
	if err != nil {
		return nil, fmt.Errorf("error creating/updating reference entity in destination: %w", err)
	}

	// 3. Get all attributes from source
	attributes, err := s.sourceRepo.FindAttributes(ctx, entityName)
	if err != nil {
		return nil, fmt.Errorf("error fetching attributes from source: %w", err)
	}

	// 4. Sync each attribute to destination
	for _, attribute := range attributes {
		attributeCode, ok := attribute["code"].(string)
		if !ok {
			return nil, fmt.Errorf("could not extract attribute code from attribute")
		}

		err := s.destRepo.SaveAttribute(ctx, entityName, attributeCode, attribute)
		if err != nil {
			return nil, fmt.Errorf("error creating/updating attribute %s in destination: %w", attributeCode, err)
		}
	}

	// 5. Get all records from source
	records, err := s.sourceRepo.FindAll(ctx, entityName)
	if err != nil {
		return nil, fmt.Errorf("error fetching records from source: %w", err)
	}

	result.TotalRecords = len(records)

	// 6. Sync each record to destination
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
