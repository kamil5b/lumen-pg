package rbac

import (
	"context"
	"fmt"
)

func (u *RBACUseCaseImplementation) GetUserAccessibleSchemas(ctx context.Context, username, database string) ([]string, error) {
	// Get accessible schemas for the user
	schemas, err := u.metadataRepo.GetAccessibleSchemas(ctx, username, database)
	if err != nil {
		return nil, fmt.Errorf("failed to get accessible schemas: %w", err)
	}

	if schemas == nil {
		return []string{}, nil
	}

	return schemas, nil
}
