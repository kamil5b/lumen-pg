package rbac

import (
	"context"
	"errors"
)

func (u *RBACUseCaseImplementation) CheckUpdatePermission(ctx context.Context, username, database, schema, table string) (bool, error) {
	return false, errors.New("not implemented")
}
