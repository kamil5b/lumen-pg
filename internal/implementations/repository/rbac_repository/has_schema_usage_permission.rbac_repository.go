package rbac_repository

import (
	"context"
	"errors"
)

func (r *RBACRepositoryImplementation) HasSchemaUsagePermission(ctx context.Context, role, database, schema string) (bool, error) {
	return false, errors.New("not implemented")
}
