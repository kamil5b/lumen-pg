package erd

import (
	"context"
	"fmt"
)

func (u *ERDUseCaseImplementation) IsSchemaEmpty(ctx context.Context, username, database, schema string) (bool, error) {
	// Get database metadata
	metadata, err := u.metadataRepo.GetMetadata(ctx, database)
	if err != nil {
		return false, fmt.Errorf("failed to get metadata: %w", err)
	}

	if metadata == nil {
		return false, fmt.Errorf("database metadata not found")
	}

	// Find the schema and check if it has any tables
	for i := range metadata.Schemas {
		if metadata.Schemas[i].Name == schema {
			return len(metadata.Schemas[i].Tables) == 0, nil
		}
	}

	return false, fmt.Errorf("schema %s not found", schema)
}
