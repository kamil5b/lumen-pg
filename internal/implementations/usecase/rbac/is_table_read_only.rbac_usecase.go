package rbac

import (
	"context"
	"fmt"
)

func (u *RBACUseCaseImplementation) IsTableReadOnly(ctx context.Context, username, database, schema, table string) (bool, error) {
	// Get table permissions for the user
	perms, err := u.metadataRepo.GetTablePermissions(ctx, username, database, schema, table)
	if err != nil {
		return false, fmt.Errorf("failed to get table permissions: %w", err)
	}

	if perms == nil {
		return false, fmt.Errorf("no permissions found for table %s", table)
	}

	// Table is read-only if user has only SELECT permission (no write permissions)
	return !perms.HasInsert && !perms.HasUpdate && !perms.HasDelete, nil
}
