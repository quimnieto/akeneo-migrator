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

// GetOptions retrieves all options for an attribute
func (r *SourceAttributeRepository) GetOptions(ctx context.Context, attributeCode string) ([]attribute.AttributeOption, error) {
	options, err := r.client.GetAttributeOptions(attributeCode)
	if err != nil {
		return nil, fmt.Errorf("error fetching options for attribute %s: %w", attributeCode, err)
	}

	// Convert from akeneo.AttributeOption to attribute.AttributeOption
	result := make([]attribute.AttributeOption, len(options))
	for i, opt := range options {
		result[i] = attribute.AttributeOption(opt)
	}

	return result, nil
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

// SaveOption creates or updates an attribute option
func (r *DestAttributeRepository) SaveOption(ctx context.Context, attributeCode, optionCode string, option attribute.AttributeOption) error {
	if err := r.client.PatchAttributeOption(attributeCode, optionCode, akeneo.AttributeOption(option)); err != nil {
		return fmt.Errorf("error saving option %s for attribute %s: %w", optionCode, attributeCode, err)
	}
	return nil
}
