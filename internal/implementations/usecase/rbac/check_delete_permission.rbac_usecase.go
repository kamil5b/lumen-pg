package rbac

import (
	"context"
	"errors"
)

func (u *RBACUseCaseImplementation) CheckDeletePermission(ctx context.Context, username, database, schema, table string) (bool, error) {
	return false, errors.New("not implemented")
}
