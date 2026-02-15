package dataview

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) GetPrimaryKeyInfo(ctx context.Context, username, database, schema, table string) ([]string, error) {
	// Check if user can access the table
	canAccess, err := u.rbacRepo.CanAccessTable(ctx, username, database, schema, table)
	if err != nil {
		return nil, err
	}
	if !canAccess {
		return nil, domain.ValidationError{Message: "cannot access table"}
	}

	// Get metadata for the database
	metadata, err := u.metadataRepo.GetMetadata(ctx, database)
	if err != nil {
		return nil, err
	}

	// Find the schema and table in metadata
	for _, schemaMetadata := range metadata.Schemas {
		if schemaMetadata.Name == schema {
			for _, tableMetadata := range schemaMetadata.Tables {
				if tableMetadata.Name == table {
					return tableMetadata.PrimaryKeys, nil
				}
			}
		}
	}

	// Table not found in metadata
	return []string{}, nil
}
