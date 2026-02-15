package setup

import (
	"context"
	"fmt"
)

func (u *SetupUseCaseImplementation) GetAllRoles(ctx context.Context) ([]string, error) {
	// Get all roles from the RBAC repository
	roles, err := u.rbacRepo.GetAllRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all roles: %w", err)
	}

	if roles == nil {
		return []string{}, nil
	}

	return roles, nil
}
