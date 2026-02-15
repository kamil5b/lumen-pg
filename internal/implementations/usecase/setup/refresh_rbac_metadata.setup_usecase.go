package setup

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *SetupUseCaseImplementation) RefreshRBACMetadata(ctx context.Context) error {
	// Invalidate all cached metadata first
	err := u.metadataRepo.InvalidateAllMetadata(ctx)
	if err != nil {
		return fmt.Errorf("failed to invalidate metadata: %w", err)
	}

	// Get all roles from the RBAC repository
	roles, err := u.rbacRepo.GetAllRoles(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all roles: %w", err)
	}

	if roles == nil || len(roles) == 0 {
		return fmt.Errorf("no roles found in database")
	}

	// Collect metadata for all roles
	rolesMetadata := make(map[string]*domain.RoleMetadata)
	for _, role := range roles {
		roleMetadata, err := u.rbacRepo.GetRoleMetadata(ctx, role)
		if err != nil {
			return fmt.Errorf("failed to get metadata for role %s: %w", role, err)
		}

		if roleMetadata != nil {
			rolesMetadata[role] = roleMetadata
		}
	}

	// Store all roles metadata in the metadata repository
	err = u.metadataRepo.StoreAllRolesMetadata(ctx, rolesMetadata)
	if err != nil {
		return fmt.Errorf("failed to store roles metadata: %w", err)
	}

	return nil
}
