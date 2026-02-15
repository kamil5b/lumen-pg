package dataview

import (
	"context"
)

func (u *DataViewUseCaseImplementation) IsTableReadOnly(ctx context.Context, username, database, schema, table string) (bool, error) {
	// Check if user has SELECT permission
	hasSelect, err := u.rbacRepo.HasSelectPermission(ctx, username, database, schema, table)
	if err != nil {
		return false, err
	}

	// If no SELECT permission, it's not readable at all
	if !hasSelect {
		return true, nil
	}

	// Check for write permissions
	hasInsert, err := u.rbacRepo.HasInsertPermission(ctx, username, database, schema, table)
	if err != nil {
		return false, err
	}
	if hasInsert {
		return false, nil
	}

	hasUpdate, err := u.rbacRepo.HasUpdatePermission(ctx, username, database, schema, table)
	if err != nil {
		return false, err
	}
	if hasUpdate {
		return false, nil
	}

	hasDelete, err := u.rbacRepo.HasDeletePermission(ctx, username, database, schema, table)
	if err != nil {
		return false, err
	}
	if hasDelete {
		return false, nil
	}

	// No write permissions found
	return true, nil
}
