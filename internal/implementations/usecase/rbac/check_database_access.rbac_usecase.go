package rbac

import (
	"context"
	"fmt"
)

func (u *RBACUseCaseImplementation) CheckDatabaseAccess(ctx context.Context, username, database string) (bool, error) {
	// Get accessible databases for the user
	accessibleDatabases, err := u.metadataRepo.GetAccessibleDatabases(ctx, username)
	if err != nil {
		return false, fmt.Errorf("failed to get accessible databases: %w", err)
	}

	// Check if the database is in the accessible list
	for _, db := range accessibleDatabases {
		if db == database {
			return true, nil
		}
	}

	return false, nil
}
