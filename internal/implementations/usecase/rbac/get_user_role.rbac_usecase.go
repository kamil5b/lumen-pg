package rbac

import (
	"context"
	"fmt"
)

func (u *RBACUseCaseImplementation) GetUserRole(ctx context.Context, username string) (string, error) {
	// Get the role of the user
	role, err := u.rbacRepo.GetUserRole(ctx, username)
	if err != nil {
		return "", fmt.Errorf("failed to get user role: %w", err)
	}

	return role, nil
}
