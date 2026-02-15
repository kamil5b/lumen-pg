package rbac

import (
	"context"
	"errors"
)

func (u *RBACUseCaseImplementation) ValidateUserAccessToResource(ctx context.Context, username, resourceType, database, schema, table string) (bool, error) {
	return false, errors.New("not implemented")
}
