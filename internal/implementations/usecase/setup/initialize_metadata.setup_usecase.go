package setup

import (
	"context"
	"fmt"
)

func (u *SetupUseCaseImplementation) InitializeMetadata(ctx context.Context, connString string) error {
	// Test the connection first
	err := u.databaseRepo.TestConnection(ctx, connString)
	if err != nil {
		return fmt.Errorf("failed to test connection: %w", err)
	}

	// Get database metadata
	metadata, err := u.databaseRepo.GetDatabaseMetadata(ctx, connString)
	if err != nil {
		return fmt.Errorf("failed to get database metadata: %w", err)
	}

	if metadata == nil {
		return fmt.Errorf("no metadata returned from database")
	}

	// Store the metadata in the metadata repository
	err = u.metadataRepo.StoreMetadata(ctx, metadata)
	if err != nil {
		return fmt.Errorf("failed to store metadata: %w", err)
	}

	return nil
}
