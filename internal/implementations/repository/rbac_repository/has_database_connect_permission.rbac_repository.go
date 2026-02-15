package rbac_repository

import (
	"context"
	"errors"
)

func (r *RBACRepositoryImplementation) HasDatabaseConnectPermission(ctx context.Context, role, database string) (bool, error) {
	return false, errors.New("not implemented")
}
