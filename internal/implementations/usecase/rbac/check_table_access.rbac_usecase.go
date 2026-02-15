package rbac

import (
	"context"
	"fmt"
)

func (u *RBACUseCaseImplementation) CheckTableAccess(ctx context.Context, username, database, schema, table string) (bool, error) {
	// Get accessible tables for the user
	accessibleTables, err := u.metadataRepo.GetAccessibleTables(ctx, username, database, schema)
	if err != nil {
		return false, fmt.Errorf("failed to get accessible tables: %w", err)
	}

	// Check if the table is in the accessible list
	for _, t := range accessibleTables {
		if t == table {
			return true, nil
		}
	}

	return false, nil
}
