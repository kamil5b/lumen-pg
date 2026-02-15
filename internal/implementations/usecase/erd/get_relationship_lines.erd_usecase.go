package erd

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *ERDUseCaseImplementation) GetRelationshipLines(ctx context.Context, username, database, schema string) ([]map[string]interface{}, error) {
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

	// Collect all relationship lines from accessible tables
	var relationshipLines []map[string]interface{}
	for _, table := range targetSchema.Tables {
		// Check RBAC permissions for this table
		canAccess, err := u.rbacRepo.CanAccessTable(ctx, username, database, schema, table.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check table access: %w", err)
		}

		if canAccess {
			// Convert foreign keys to relationship line format
			for _, fk := range table.ForeignKeys {
				line := map[string]interface{}{
					"sourceTable":    table.Name,
					"sourceColumn":   fk.ColumnName,
					"targetTable":    fk.ReferencedTable,
					"targetColumn":   fk.ReferencedColumn,
					"targetSchema":   fk.ReferencedSchema,
					"targetDatabase": fk.ReferencedDatabase,
				}
				relationshipLines = append(relationshipLines, line)
			}
		}
	}

	return relationshipLines, nil
}
