package akeneo

import (
	"context"
	"fmt"

	"akeneo-migrator/internal/category"
	"akeneo-migrator/internal/platform/client/akeneo"
)

// SourceCategoryRepository implements category.SourceRepository for Akeneo
type SourceCategoryRepository struct {
	client *akeneo.Client
}

// NewSourceCategoryRepository creates a new source category repository
func NewSourceCategoryRepository(client *akeneo.Client) category.SourceRepository {
	return &SourceCategoryRepository{
		client: client,
	}
}

// FindByCode retrieves a category by its code
func (r *SourceCategoryRepository) FindByCode(ctx context.Context, code string) (category.Category, error) {
	cat, err := r.client.GetCategory(code)
	if err != nil {
		return nil, fmt.Errorf("error fetching category %s: %w", code, err)
	}
	return category.Category(cat), nil
}

// DestCategoryRepository implements category.DestRepository for Akeneo
type DestCategoryRepository struct {
	client *akeneo.Client
}

// NewDestCategoryRepository creates a new destination category repository
func NewDestCategoryRepository(client *akeneo.Client) category.DestRepository {
	return &DestCategoryRepository{
		client: client,
	}
}

// Save creates or updates a category
func (r *DestCategoryRepository) Save(ctx context.Context, code string, cat category.Category) error {
	if err := r.client.PatchCategory(code, akeneo.Category(cat)); err != nil {
		return fmt.Errorf("error saving category %s: %w", code, err)
	}
	return nil
}
