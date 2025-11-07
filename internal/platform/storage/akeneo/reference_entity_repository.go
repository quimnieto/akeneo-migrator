package akeneo

import (
	"context"

	"akeneo-migrator/internal/platform/client/akeneo"
	"akeneo-migrator/internal/reference_entity"
)

// SourceReferenceEntityRepository implements the read-only repository for the source
type SourceReferenceEntityRepository struct {
	client *akeneo.Client
}

// NewSourceReferenceEntityRepository creates a new instance of the source repository
func NewSourceReferenceEntityRepository(client *akeneo.Client) *SourceReferenceEntityRepository {
	return &SourceReferenceEntityRepository{
		client: client,
	}
}

// FindEntity retrieves a Reference Entity definition
func (r *SourceReferenceEntityRepository) FindEntity(ctx context.Context, entityCode string) (reference_entity.Entity, error) {
	entity, err := r.client.GetReferenceEntity(entityCode)
	if err != nil {
		return nil, err
	}

	return reference_entity.Entity(entity), nil
}

// FindAttributes retrieves all attributes from a Reference Entity
func (r *SourceReferenceEntityRepository) FindAttributes(ctx context.Context, entityCode string) ([]reference_entity.Attribute, error) {
	attributes, err := r.client.GetReferenceEntityAttributes(entityCode)
	if err != nil {
		return nil, err
	}

	// Convert from akeneo.ReferenceEntityAttribute to reference_entity.Attribute
	result := make([]reference_entity.Attribute, len(attributes))
	for i, attr := range attributes {
		result[i] = reference_entity.Attribute(attr)
	}

	return result, nil
}

// FindAll retrieves all records from a Reference Entity
func (r *SourceReferenceEntityRepository) FindAll(ctx context.Context, entityName string) ([]reference_entity.Record, error) {
	records, err := r.client.GetReferenceEntityRecords(entityName)
	if err != nil {
		return nil, err
	}

	// Convert from akeneo.ReferenceEntityRecord to reference_entity.Record
	result := make([]reference_entity.Record, len(records))
	for i, record := range records {
		result[i] = reference_entity.Record(record)
	}

	return result, nil
}

// DestReferenceEntityRepository implements the read/write repository for the destination
type DestReferenceEntityRepository struct {
	client *akeneo.Client
}

// NewDestReferenceEntityRepository creates a new instance of the destination repository
func NewDestReferenceEntityRepository(client *akeneo.Client) *DestReferenceEntityRepository {
	return &DestReferenceEntityRepository{
		client: client,
	}
}

// FindEntity retrieves a Reference Entity definition
func (r *DestReferenceEntityRepository) FindEntity(ctx context.Context, entityCode string) (reference_entity.Entity, error) {
	entity, err := r.client.GetReferenceEntity(entityCode)
	if err != nil {
		return nil, err
	}

	return reference_entity.Entity(entity), nil
}

// SaveEntity creates or updates a Reference Entity definition
func (r *DestReferenceEntityRepository) SaveEntity(ctx context.Context, entityCode string, entity reference_entity.Entity) error {
	// Convert from reference_entity.Entity to akeneo.ReferenceEntity
	akeneoEntity := akeneo.ReferenceEntity(entity)
	return r.client.PatchReferenceEntity(entityCode, akeneoEntity)
}

// FindAttributes retrieves all attributes from a Reference Entity
func (r *DestReferenceEntityRepository) FindAttributes(ctx context.Context, entityCode string) ([]reference_entity.Attribute, error) {
	attributes, err := r.client.GetReferenceEntityAttributes(entityCode)
	if err != nil {
		return nil, err
	}

	// Convert from akeneo.ReferenceEntityAttribute to reference_entity.Attribute
	result := make([]reference_entity.Attribute, len(attributes))
	for i, attr := range attributes {
		result[i] = reference_entity.Attribute(attr)
	}

	return result, nil
}

// SaveAttribute creates or updates a Reference Entity attribute
func (r *DestReferenceEntityRepository) SaveAttribute(ctx context.Context, entityCode string, attributeCode string, attribute reference_entity.Attribute) error {
	// Convert from reference_entity.Attribute to akeneo.ReferenceEntityAttribute
	akeneoAttribute := akeneo.ReferenceEntityAttribute(attribute)
	return r.client.PatchReferenceEntityAttribute(entityCode, attributeCode, akeneoAttribute)
}

// FindAll retrieves all records from a Reference Entity
func (r *DestReferenceEntityRepository) FindAll(ctx context.Context, entityName string) ([]reference_entity.Record, error) {
	records, err := r.client.GetReferenceEntityRecords(entityName)
	if err != nil {
		return nil, err
	}

	// Convert from akeneo.ReferenceEntityRecord to reference_entity.Record
	result := make([]reference_entity.Record, len(records))
	for i, record := range records {
		result[i] = reference_entity.Record(record)
	}

	return result, nil
}

// Save creates or updates a record in a Reference Entity
func (r *DestReferenceEntityRepository) Save(ctx context.Context, entityName string, code string, record reference_entity.Record) error {
	// Convert from reference_entity.Record to akeneo.ReferenceEntityRecord
	akeneoRecord := akeneo.ReferenceEntityRecord(record)
	return r.client.PatchReferenceEntityRecord(entityName, code, akeneoRecord)
}
