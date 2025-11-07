package syncing_since

import (
	"context"
	"fmt"

	"akeneo-migrator/internal/product"
	"akeneo-migrator/internal/product/syncing"
)

// Service handles the synchronization of updated products
type Service struct {
	sourceRepo     product.SourceRepository
	destRepo       product.DestRepository
	syncingService *syncing.Service
}

// NewService creates a new instance of the sync since service
func NewService(sourceRepo product.SourceRepository, destRepo product.DestRepository) *Service {
	return &Service{
		sourceRepo:     sourceRepo,
		destRepo:       destRepo,
		syncingService: syncing.NewService(sourceRepo, destRepo),
	}
}

// SyncResult contains the result of syncing updated products
type SyncResult struct {
	UpdatedSince   string
	ProductsSynced int
	ModelsSynced   int
	TotalSynced    int
	Errors         []string
	Success        bool
}

// Sync synchronizes all products and models updated since a specific date
// Optimized: Only syncs root/common products (without parent), the hierarchy sync handles the rest
func (s *Service) Sync(ctx context.Context, updatedSince string) (*SyncResult, error) {
	result := &SyncResult{
		UpdatedSince: updatedSince,
		Success:      true,
	}

	fmt.Printf("📅 Fetching products updated since: %s\n", updatedSince)

	// 1. Get all updated product models
	models, err := s.sourceRepo.FindModelsUpdatedSince(ctx, updatedSince)
	if err != nil {
		return nil, fmt.Errorf("error fetching updated models: %w", err)
	}

	fmt.Printf("   📦 Found %d updated models\n", len(models))

	// 2. Get all updated products
	products, err := s.sourceRepo.FindProductsUpdatedSince(ctx, updatedSince)
	if err != nil {
		return nil, fmt.Errorf("error fetching updated products: %w", err)
	}

	fmt.Printf("   📦 Found %d updated products\n", len(products))

	// 3. Filter to get only root/common models (without parent)
	commonModels := s.filterCommonModels(models)
	fmt.Printf("   🎯 Identified %d common models (root level)\n", len(commonModels))

	// 4. Filter to get only root/common products (without parent)
	commonProducts := s.filterCommonProducts(products)
	fmt.Printf("   🎯 Identified %d common products (root level)\n", len(commonProducts))

	// 5. Track synced hierarchies to avoid duplicates
	syncedHierarchies := make(map[string]bool)

	// 6. Sync each common model hierarchy
	for _, model := range commonModels {
		code, ok := model["code"].(string)
		if !ok {
			result.Errors = append(result.Errors, "could not extract model code")
			continue
		}

		if syncedHierarchies[code] {
			fmt.Printf("   ⏭️  Skipping already synced hierarchy: %s\n", code)
			continue
		}

		fmt.Printf("   🔄 Syncing model hierarchy: %s\n", code)

		hierarchyResult, err := s.syncingService.Sync(ctx, code)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("error syncing model %s: %v", code, err))
			continue
		}

		syncedHierarchies[code] = true
		result.ModelsSynced += hierarchyResult.ModelsSynced
		result.ProductsSynced += hierarchyResult.ProductsSynced
	}

	// 7. Sync each common product hierarchy
	for _, prod := range commonProducts {
		identifier, ok := prod["identifier"].(string)
		if !ok {
			result.Errors = append(result.Errors, "could not extract product identifier")
			continue
		}

		if syncedHierarchies[identifier] {
			fmt.Printf("   ⏭️  Skipping already synced hierarchy: %s\n", identifier)
			continue
		}

		fmt.Printf("   🔄 Syncing product hierarchy: %s\n", identifier)

		hierarchyResult, err := s.syncingService.Sync(ctx, identifier)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("error syncing product %s: %v", identifier, err))
			continue
		}

		syncedHierarchies[identifier] = true
		result.ModelsSynced += hierarchyResult.ModelsSynced
		result.ProductsSynced += hierarchyResult.ProductsSynced
	}

	result.TotalSynced = result.ModelsSynced + result.ProductsSynced

	if len(result.Errors) > 0 {
		result.Success = false
	}

	return result, nil
}

// filterCommonModels returns only models without parent (root/common models)
func (s *Service) filterCommonModels(models []product.ProductModel) []product.ProductModel {
	var commonModels []product.ProductModel
	for _, model := range models {
		parent, hasParent := model["parent"].(string)
		// Include if no parent or parent is empty/null
		if !hasParent || parent == "" || parent == "null" {
			commonModels = append(commonModels, model)
		}
	}
	return commonModels
}

// filterCommonProducts returns only products without parent (root/common products)
func (s *Service) filterCommonProducts(products []product.Product) []product.Product {
	var commonProducts []product.Product
	for _, prod := range products {
		parent, hasParent := prod["parent"].(string)
		// Include if no parent or parent is empty/null
		if !hasParent || parent == "" || parent == "null" {
			commonProducts = append(commonProducts, prod)
		}
	}
	return commonProducts
}
