package setup

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *SetupUseCaseImplementation) GetRoleAccessibility(ctx context.Context, role string) (*domain.RoleMetadata, error) {
	// Get role metadata from the RBAC repository
	metadata, err := u.rbacRepo.GetRoleMetadata(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("failed to get role metadata: %w", err)
	}

	if metadata == nil {
		return nil, fmt.Errorf("no metadata found for role %s", role)
	}

	return metadata, nil
}
