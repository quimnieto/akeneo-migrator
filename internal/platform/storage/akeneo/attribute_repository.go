package akeneo

import (
	"context"
	"fmt"

	"akeneo-migrator/internal/attribute"
	"akeneo-migrator/internal/platform/client/akeneo"
)

// SourceAttributeRepository implements attribute.SourceRepository for Akeneo
type SourceAttributeRepository struct {
	client *akeneo.Client
}

// NewSourceAttributeRepository creates a new source attribute repository
func NewSourceAttributeRepository(client *akeneo.Client) attribute.SourceRepository {
	return &SourceAttributeRepository{
		client: client,
	}
}

// FindByCode retrieves an attribute by its code
func (r *SourceAttributeRepository) FindByCode(ctx context.Context, code string) (attribute.Attribute, error) {
	attr, err := r.client.GetAttribute(code)
	if err != nil {
		return nil, fmt.Errorf("error fetching attribute %s: %w", code, err)
	}
	return attribute.Attribute(attr), nil
}

// DestAttributeRepository implements attribute.DestRepository for Akeneo
type DestAttributeRepository struct {
	client *akeneo.Client
}

// NewDestAttributeRepository creates a new destination attribute repository
func NewDestAttributeRepository(client *akeneo.Client) attribute.DestRepository {
	return &DestAttributeRepository{
		client: client,
	}
}

// Save creates or updates an attribute
func (r *DestAttributeRepository) Save(ctx context.Context, code string, attr attribute.Attribute) error {
	if err := r.client.PatchAttribute(code, akeneo.Attribute(attr)); err != nil {
		return fmt.Errorf("error saving attribute %s: %w", code, err)
	}
	return nil
}
