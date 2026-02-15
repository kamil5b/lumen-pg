package rbac

import (
	"context"
	"fmt"
)

func (u *RBACUseCaseImplementation) GetUserAccessibleDatabases(ctx context.Context, username string) ([]string, error) {
	// Get accessible databases for the user
	databases, err := u.metadataRepo.GetAccessibleDatabases(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get accessible databases: %w", err)
	}

	if databases == nil {
		return []string{}, nil
	}

	return databases, nil
}
