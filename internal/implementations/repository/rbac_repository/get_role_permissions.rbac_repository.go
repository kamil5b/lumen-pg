package rbac_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (r *RBACRepositoryImplementation) GetRolePermissions(ctx context.Context, role, database, schema, table string) (*domain.PermissionSet, error) {
	return nil, errors.New("not implemented")
}
