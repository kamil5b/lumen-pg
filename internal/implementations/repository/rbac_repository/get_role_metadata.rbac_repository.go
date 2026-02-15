package rbac_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (r *RBACRepositoryImplementation) GetRoleMetadata(ctx context.Context, role string) (*domain.RoleMetadata, error) {
	return nil, errors.New("not implemented")
}
