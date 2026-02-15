package rbac_repository

import (
	"context"
	"errors"
)

func (r *RBACRepositoryImplementation) GetUserRole(ctx context.Context, username string) (string, error) {
	return "", errors.New("not implemented")
}
