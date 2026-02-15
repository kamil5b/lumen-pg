package erd

import (
	"context"
	"fmt"
)

func (u *ERDUseCaseImplementation) GetAvailableSchemas(ctx context.Context, username, database string) ([]string, error) {
	// Get database metadata
	metadata, err := u.metadataRepo.GetMetadata(ctx, database)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	if metadata == nil {
		return nil, fmt.Errorf("database metadata not found")
	}

	// Collect all schemas where user has usage permission
	var availableSchemas []string
	for _, schemaMetadata := range metadata.Schemas {
		hasPermission, err := u.rbacRepo.HasSchemaUsagePermission(ctx, username, database, schemaMetadata.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check schema permission: %w", err)
		}

		if hasPermission {
			availableSchemas = append(availableSchemas, schemaMetadata.Name)
		}
	}

	return availableSchemas, nil
}
