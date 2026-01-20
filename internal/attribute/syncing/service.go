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
	Code          string
	Success       bool
	Error         string
	OptionsSynced int
	OptionsErrors []string
}

// Sync synchronizes a single attribute from source to destination
func (s *Service) Sync(ctx context.Context, code string) (*SyncResult, error) {
	result := &SyncResult{
		Code:          code,
		Success:       false,
		OptionsSynced: 0,
		OptionsErrors: []string{},
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

	// 3. Check if attribute is of type select (simple or multi)
	attributeType, ok := attributeData["type"].(string)
	if ok && (attributeType == "pim_catalog_simpleselect" || attributeType == "pim_catalog_multiselect") {
		// 4. Get attribute options from source
		options, err := s.sourceRepo.GetOptions(ctx, code)
		if err != nil {
			// Log error but don't fail the entire sync
			result.OptionsErrors = append(result.OptionsErrors, fmt.Sprintf("error fetching options: %v", err))
		} else {
			// 5. Sync each option to destination
			for _, option := range options {
				optionCode, ok := option["code"].(string)
				if !ok {
					result.OptionsErrors = append(result.OptionsErrors, "option without code field")
					continue
				}

				err := s.destRepo.SaveOption(ctx, code, optionCode, option)
				if err != nil {
					result.OptionsErrors = append(result.OptionsErrors, fmt.Sprintf("option %s: %v", optionCode, err))
				} else {
					result.OptionsSynced++
				}
			}
		}
	}

	result.Success = true
	return result, nil
}
