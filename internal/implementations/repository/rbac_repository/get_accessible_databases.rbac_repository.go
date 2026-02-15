package rbac_repository

import (
	"context"
	"errors"
)

func (r *RBACRepositoryImplementation) GetAccessibleDatabases(ctx context.Context, role string) ([]string, error) {
	return nil, errors.New("not implemented")
}
