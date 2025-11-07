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
