package rbac_repository

import (
	"context"
	"errors"
)

func (r *RBACRepositoryImplementation) GetAccessibleSchemas(ctx context.Context, role, database string) ([]string, error) {
	return nil, errors.New("not implemented")
}
