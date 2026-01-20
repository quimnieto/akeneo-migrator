package syncing

import (
	"context"
	"fmt"

	"akeneo-migrator/internal/family"
)

// Service handles family synchronization
type Service struct {
	sourceRepo family.SourceRepository
	destRepo   family.DestRepository
}

// NewService creates a new family sync service
func NewService(sourceRepo family.SourceRepository, destRepo family.DestRepository) *Service {
	return &Service{
		sourceRepo: sourceRepo,
		destRepo:   destRepo,
	}
}

// SyncResult contains the result of a sync operation
type SyncResult struct {
	Code           string
	Success        bool
	Error          string
	VariantsSynced int
	VariantsErrors []string
}

// Sync synchronizes a single family from source to destination
func (s *Service) Sync(ctx context.Context, code string) (*SyncResult, error) {
	result := &SyncResult{
		Code:           code,
		Success:        false,
		VariantsSynced: 0,
		VariantsErrors: []string{},
	}

	// 1. Get family from source
	familyData, err := s.sourceRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error fetching family from source: %w", err)
	}

	// 2. Save family to destination
	err = s.destRepo.Save(ctx, code, familyData)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, fmt.Errorf("error saving family to destination: %w", err)
	}

	// 3. Get family variants from source
	variants, err := s.sourceRepo.GetVariants(ctx, code)
	if err != nil {
		// Log error but don't fail the entire sync
		result.VariantsErrors = append(result.VariantsErrors, fmt.Sprintf("error fetching variants: %v", err))
	} else {
		// 4. Sync each variant to destination
		for _, variant := range variants {
			variantCode, ok := variant["code"].(string)
			if !ok {
				result.VariantsErrors = append(result.VariantsErrors, "variant without code field")
				continue
			}

			err := s.destRepo.SaveVariant(ctx, code, variantCode, variant)
			if err != nil {
				result.VariantsErrors = append(result.VariantsErrors, fmt.Sprintf("variant %s: %v", variantCode, err))
			} else {
				result.VariantsSynced++
			}
		}
	}

	result.Success = true
	return result, nil
}
