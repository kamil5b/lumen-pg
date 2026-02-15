package erd

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *ERDUseCaseImplementation) GetERDTables(ctx context.Context, username, database, schema string) ([]domain.TableMetadata, error) {
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

	// Filter tables based on RBAC permissions
	var accessibleTables []domain.TableMetadata
	for _, table := range targetSchema.Tables {
		canAccess, err := u.rbacRepo.CanAccessTable(ctx, username, database, schema, table.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check table access: %w", err)
		}

		if canAccess {
			accessibleTables = append(accessibleTables, table)
		}
	}

	return accessibleTables, nil
}
