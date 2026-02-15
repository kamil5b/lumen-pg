package rbac

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *RBACUseCaseImplementation) VerifyUserPermissions(ctx context.Context, username, database, schema, table string) (*domain.PermissionSet, error) {
	// Get table permissions for the user
	perms, err := u.metadataRepo.GetTablePermissions(ctx, username, database, schema, table)
	if err != nil {
		return nil, fmt.Errorf("failed to get table permissions: %w", err)
	}

	if perms == nil {
		return nil, fmt.Errorf("no permissions found for table %s", table)
	}

	// Convert AccessibleTable to PermissionSet
	permissionSet := &domain.PermissionSet{
		CanSelect: perms.HasSelect,
		CanInsert: perms.HasInsert,
		CanUpdate: perms.HasUpdate,
		CanDelete: perms.HasDelete,
	}

	return permissionSet, nil
}
