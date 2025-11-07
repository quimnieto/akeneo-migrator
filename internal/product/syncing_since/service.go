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
// Memory-efficient: Processes products/models in batches using streaming
// Logic: For each updated product/model, finds its root and syncs the entire hierarchy
func (s *Service) Sync(ctx context.Context, updatedSince string) (*SyncResult, error) {
	result := &SyncResult{
		UpdatedSince: updatedSince,
		Success:      true,
	}

	fmt.Printf("📅 Syncing products updated since: %s (streaming mode)\n", updatedSince)

	// Track synced hierarchies to avoid duplicates
	syncedHierarchies := make(map[string]bool)

	batchSize := 100 // Process 100 items at a time
	modelsProcessed := 0
	productsProcessed := 0

	// 1. Stream and process product models in batches
	fmt.Println("   📦 Processing product models...")
	err := s.sourceRepo.StreamModelsUpdatedSince(ctx, updatedSince, batchSize, func(models []product.ProductModel) error {
		for _, model := range models {
			code, ok := model["code"].(string)
			if !ok {
				result.Errors = append(result.Errors, "could not extract model code")
				continue
			}

			// Find the root of this model's hierarchy
			root := s.findModelRoot(ctx, model)

			// Skip if already synced
			if syncedHierarchies[root] {
				continue
			}

			fmt.Printf("   🔄 Syncing hierarchy from root: %s (triggered by model: %s)\n", root, code)

			hierarchyResult, syncErr := s.syncingService.Sync(ctx, root)
			if syncErr != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("error syncing root %s: %v", root, syncErr))
				continue
			}

			syncedHierarchies[root] = true
			result.ModelsSynced += hierarchyResult.ModelsSynced
			result.ProductsSynced += hierarchyResult.ProductsSynced
			modelsProcessed++
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error streaming updated models: %w", err)
	}

	fmt.Printf("   ✅ Processed %d models (found their roots)\n", modelsProcessed)

	// 2. Stream and process products in batches
	fmt.Println("   📦 Processing products...")
	err = s.sourceRepo.StreamProductsUpdatedSince(ctx, updatedSince, batchSize, func(products []product.Product) error {
		for _, prod := range products {
			identifier, ok := prod["identifier"].(string)
			if !ok {
				result.Errors = append(result.Errors, "could not extract product identifier")
				continue
			}

			// Find the root of this product's hierarchy
			root := s.findProductRoot(ctx, prod)

			// Skip if already synced
			if syncedHierarchies[root] {
				continue
			}

			fmt.Printf("   🔄 Syncing hierarchy from root: %s (triggered by product: %s)\n", root, identifier)

			hierarchyResult, syncErr := s.syncingService.Sync(ctx, root)
			if syncErr != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("error syncing root %s: %v", root, syncErr))
				continue
			}

			syncedHierarchies[root] = true
			result.ModelsSynced += hierarchyResult.ModelsSynced
			result.ProductsSynced += hierarchyResult.ProductsSynced
			productsProcessed++
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error streaming updated products: %w", err)
	}

	fmt.Printf("   ✅ Processed %d products (found their roots)\n", productsProcessed)

	result.TotalSynced = result.ModelsSynced + result.ProductsSynced

	if len(result.Errors) > 0 {
		result.Success = false
	}

	return result, nil
}

// findModelRoot navigates up the hierarchy to find the root model
func (s *Service) findModelRoot(ctx context.Context, model product.ProductModel) string {
	code, _ := model["code"].(string)
	parent, hasParent := model["parent"].(string)

	// If no parent, this is the root
	if !hasParent || parent == "" || parent == "null" {
		return code
	}

	// Navigate up to find the root
	parentModel, err := s.sourceRepo.FindModelByCode(ctx, parent)
	if err != nil {
		// If we can't find the parent, treat current as root
		return code
	}

	// Recursively find the root
	return s.findModelRoot(ctx, parentModel)
}

// findProductRoot navigates up the hierarchy to find the root (model or product)
func (s *Service) findProductRoot(ctx context.Context, prod product.Product) string {
	identifier, _ := prod["identifier"].(string)
	parent, hasParent := prod["parent"].(string)

	// If no parent, this is the root
	if !hasParent || parent == "" || parent == "null" {
		return identifier
	}

	// Try to find parent as a model first
	parentModel, err := s.sourceRepo.FindModelByCode(ctx, parent)
	if err == nil {
		// Found as model, navigate up from there
		return s.findModelRoot(ctx, parentModel)
	}

	// Try as product
	parentProduct, err := s.sourceRepo.FindByIdentifier(ctx, parent)
	if err != nil {
		// If we can't find the parent, treat current as root
		return identifier
	}

	// Recursively find the root
	return s.findProductRoot(ctx, parentProduct)
}
