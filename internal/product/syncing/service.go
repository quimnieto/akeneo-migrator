package syncing

import (
	"context"
	"fmt"

	"akeneo-migrator/internal/product"
)

// Service handles the synchronization logic for Products
type Service struct {
	sourceRepo product.SourceRepository
	destRepo   product.DestRepository
}

// NewService creates a new instance of the synchronization service
func NewService(sourceRepo product.SourceRepository, destRepo product.DestRepository) *Service {
	return &Service{
		sourceRepo: sourceRepo,
		destRepo:   destRepo,
	}
}

// SyncResult contains the result of a synchronization operation
type SyncResult struct {
	Identifier     string
	Success        bool
	Error          string
	ModelsSynced   int
	ProductsSynced int
	TotalSynced    int
}

// Sync synchronizes a complete product hierarchy (common → models → products)
func (s *Service) Sync(ctx context.Context, commonIdentifier string) (*SyncResult, error) {
	result := &SyncResult{
		Identifier: commonIdentifier,
	}

	// 1. Sync the common product/model
	fmt.Printf("   📦 Syncing common: %s\n", commonIdentifier)

	// Try as product first
	commonProduct, err := s.sourceRepo.FindByIdentifier(ctx, commonIdentifier)
	if err == nil {
		// It's a product (simple type)
		if err := s.destRepo.Save(ctx, commonIdentifier, commonProduct); err != nil {
			return nil, fmt.Errorf("error saving common product: %w", err)
		}
		result.ProductsSynced++

		// Sync child products
		if err := s.syncChildProducts(ctx, commonIdentifier, result); err != nil {
			return nil, err
		}
	} else {
		// Try as product model (configurable type)
		commonModel, modelErr := s.sourceRepo.FindModelByCode(ctx, commonIdentifier)
		if modelErr != nil {
			return nil, fmt.Errorf("common '%s' not found as product or model: %w", commonIdentifier, modelErr)
		}

		if err := s.destRepo.SaveModel(ctx, commonIdentifier, commonModel); err != nil {
			return nil, fmt.Errorf("error saving common model: %w", err)
		}
		result.ModelsSynced++

		// Sync child models
		if err := s.syncChildModels(ctx, commonIdentifier, result); err != nil {
			return nil, err
		}

		// Sync all variant products under all models
		if err := s.syncVariantProducts(ctx, commonIdentifier, result); err != nil {
			return nil, err
		}
	}

	result.TotalSynced = result.ModelsSynced + result.ProductsSynced
	result.Success = true
	return result, nil
}

// syncChildProducts syncs all products that have the given parent
func (s *Service) syncChildProducts(ctx context.Context, parentCode string, result *SyncResult) error {
	products, err := s.sourceRepo.FindProductsByParent(ctx, parentCode)
	if err != nil {
		return fmt.Errorf("error fetching child products: %w", err)
	}

	fmt.Printf("   👶 Found %d child products\n", len(products))

	for _, prod := range products {
		identifier, _ := prod["identifier"].(string)
		if identifier == "" {
			continue
		}

		if err := s.destRepo.Save(ctx, identifier, prod); err != nil {
			fmt.Printf("   ⚠️  Error syncing product %s: %v\n", identifier, err)
			continue
		}

		fmt.Printf("   ✅ Synced product: %s\n", identifier)
		result.ProductsSynced++
	}

	return nil
}

// syncChildModels syncs all product models that have the given parent
func (s *Service) syncChildModels(ctx context.Context, parentCode string, result *SyncResult) error {
	models, err := s.sourceRepo.FindModelsByParent(ctx, parentCode)
	if err != nil {
		return fmt.Errorf("error fetching child models: %w", err)
	}

	fmt.Printf("   📋 Found %d child models\n", len(models))

	for _, model := range models {
		code, _ := model["code"].(string)
		if code == "" {
			continue
		}

		if err := s.destRepo.SaveModel(ctx, code, model); err != nil {
			fmt.Printf("   ⚠️  Error syncing model %s: %v\n", code, err)
			continue
		}

		fmt.Printf("   ✅ Synced model: %s\n", code)
		result.ModelsSynced++
	}

	return nil
}

// syncVariantProducts syncs all variant products under all models of a common
func (s *Service) syncVariantProducts(ctx context.Context, commonCode string, result *SyncResult) error {
	// Get all models under the common
	models, err := s.sourceRepo.FindModelsByParent(ctx, commonCode)
	if err != nil {
		return fmt.Errorf("error fetching models for variants: %w", err)
	}

	// For each model, get its variant products
	for _, model := range models {
		modelCode, _ := model["code"].(string)
		if modelCode == "" {
			continue
		}

		products, err := s.sourceRepo.FindProductsByParent(ctx, modelCode)
		if err != nil {
			fmt.Printf("   ⚠️  Error fetching variants for model %s: %v\n", modelCode, err)
			continue
		}

		fmt.Printf("   🔸 Found %d variants for model %s\n", len(products), modelCode)

		for _, prod := range products {
			identifier, _ := prod["identifier"].(string)
			if identifier == "" {
				continue
			}

			if err := s.destRepo.Save(ctx, identifier, prod); err != nil {
				fmt.Printf("   ⚠️  Error syncing variant %s: %v\n", identifier, err)
				continue
			}

			fmt.Printf("   ✅ Synced variant: %s\n", identifier)
			result.ProductsSynced++
		}
	}

	return nil
}
