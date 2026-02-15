package erd

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *ERDUseCaseImplementation) GetERDRelationships(ctx context.Context, username, database, schema string) ([]domain.ForeignKeyMetadata, error) {
	// Get database metadata
	metadata, err := u.metadataRepo.GetMetadata(ctx, database)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	if metadata == nil {
		return nil, fmt.Errorf("database metadata not found")
	}

	// Find the schema
	var targetSchema *domain.SchemaMetadata
	for i := range metadata.Schemas {
		if metadata.Schemas[i].Name == schema {
			targetSchema = &metadata.Schemas[i]
			break
		}
	}

	if targetSchema == nil {
		return nil, fmt.Errorf("schema %s not found", schema)
	}

	// Collect all foreign key relationships from accessible tables
	var relationships []domain.ForeignKeyMetadata
	for _, table := range targetSchema.Tables {
		// Check RBAC permissions for this table
		canAccess, err := u.rbacRepo.CanAccessTable(ctx, username, database, schema, table.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check table access: %w", err)
		}

		if canAccess {
			relationships = append(relationships, table.ForeignKeys...)
		}
	}

	return relationships, nil
}
