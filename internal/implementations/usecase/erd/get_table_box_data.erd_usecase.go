package erd

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *ERDUseCaseImplementation) GetTableBoxData(ctx context.Context, username, database, schema, table string) (*domain.TableMetadata, error) {
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

	// Find the table
	var targetTable *domain.TableMetadata
	for i := range targetSchema.Tables {
		if targetSchema.Tables[i].Name == table {
			targetTable = &targetSchema.Tables[i]
			break
		}
	}

	if targetTable == nil {
		return nil, fmt.Errorf("table %s not found in schema %s", table, schema)
	}

	// Check RBAC permissions
	canAccess, err := u.rbacRepo.CanAccessTable(ctx, username, database, schema, table)
	if err != nil {
		return nil, fmt.Errorf("failed to check table access: %w", err)
	}

	if !canAccess {
		return nil, fmt.Errorf("access denied to table %s", table)
	}

	return targetTable, nil
}
