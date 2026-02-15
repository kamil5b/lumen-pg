package setup

import (
	"context"
	"fmt"
)

func (u *SetupUseCaseImplementation) RefreshMetadata(ctx context.Context) error {
	// Invalidate all cached metadata first
	err := u.metadataRepo.InvalidateAllMetadata(ctx)
	if err != nil {
		return fmt.Errorf("failed to invalidate metadata: %w", err)
	}

	// Get fresh database metadata
	// We need a connection string, but it's not passed in this method
	// We'll try to get metadata with empty connection string and let the repository handle it
	metadata, err := u.databaseRepo.GetDatabaseMetadata(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get database metadata: %w", err)
	}

	if metadata == nil {
		return fmt.Errorf("no metadata returned from database")
	}

	// Store the refreshed metadata
	err = u.metadataRepo.StoreMetadata(ctx, metadata)
	if err != nil {
		return fmt.Errorf("failed to store refreshed metadata: %w", err)
	}

	return nil
}
