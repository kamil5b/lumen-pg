package rbac_repository

import (
	"context"
	"errors"
)

func (r *RBACRepositoryImplementation) IsReadOnlyRole(ctx context.Context, role, database, schema, table string) (bool, error) {
	return false, errors.New("not implemented")
}
