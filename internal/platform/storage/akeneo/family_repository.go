package akeneo

import (
	"context"
	"fmt"

	"akeneo-migrator/internal/family"
	"akeneo-migrator/internal/platform/client/akeneo"
)

// SourceFamilyRepository implements family.SourceRepository for Akeneo
type SourceFamilyRepository struct {
	client *akeneo.Client
}

// NewSourceFamilyRepository creates a new source family repository
func NewSourceFamilyRepository(client *akeneo.Client) family.SourceRepository {
	return &SourceFamilyRepository{
		client: client,
	}
}

// FindByCode retrieves a family by its code
func (r *SourceFamilyRepository) FindByCode(ctx context.Context, code string) (family.Family, error) {
	fam, err := r.client.GetFamily(code)
	if err != nil {
		return nil, fmt.Errorf("error fetching family %s: %w", code, err)
	}
	return family.Family(fam), nil
}

// GetVariants retrieves all variants for a family
func (r *SourceFamilyRepository) GetVariants(ctx context.Context, familyCode string) ([]family.FamilyVariant, error) {
	variants, err := r.client.GetFamilyVariants(familyCode)
	if err != nil {
		return nil, fmt.Errorf("error fetching variants for family %s: %w", familyCode, err)
	}

	// Convert from akeneo.FamilyVariant to family.FamilyVariant
	result := make([]family.FamilyVariant, len(variants))
	for i, variant := range variants {
		result[i] = family.FamilyVariant(variant)
	}

	return result, nil
}

// DestFamilyRepository implements family.DestRepository for Akeneo
type DestFamilyRepository struct {
	client *akeneo.Client
}

// NewDestFamilyRepository creates a new destination family repository
func NewDestFamilyRepository(client *akeneo.Client) family.DestRepository {
	return &DestFamilyRepository{
		client: client,
	}
}

// Save creates or updates a family
func (r *DestFamilyRepository) Save(ctx context.Context, code string, fam family.Family) error {
	if err := r.client.PatchFamily(code, akeneo.Family(fam)); err != nil {
		return fmt.Errorf("error saving family %s: %w", code, err)
	}
	return nil
}

// SaveVariant creates or updates a family variant
func (r *DestFamilyRepository) SaveVariant(ctx context.Context, familyCode, variantCode string, variant family.FamilyVariant) error {
	if err := r.client.PatchFamilyVariant(familyCode, variantCode, akeneo.FamilyVariant(variant)); err != nil {
		return fmt.Errorf("error saving variant %s for family %s: %w", variantCode, familyCode, err)
	}
	return nil
}
