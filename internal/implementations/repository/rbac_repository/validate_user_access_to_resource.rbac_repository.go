package rbac_repository

import (
	"context"
	"errors"
)

func (r *RBACRepositoryImplementation) ValidateUserAccessToResource(ctx context.Context, username, resourceType, database, schema, table string) (bool, error) {
	return false, errors.New("not implemented")
}
