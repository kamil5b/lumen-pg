package rbac

import (
	"context"
	"fmt"
)

func (u *RBACUseCaseImplementation) CheckSchemaAccess(ctx context.Context, username, database, schema string) (bool, error) {
	// Get accessible schemas for the user
	accessibleSchemas, err := u.metadataRepo.GetAccessibleSchemas(ctx, username, database)
	if err != nil {
		return false, fmt.Errorf("failed to get accessible schemas: %w", err)
	}

	// Check if the schema is in the accessible list
	for _, s := range accessibleSchemas {
		if s == schema {
			return true, nil
		}
	}

	return false, nil
}
