package rbac

import (
	"context"
	"errors"
)

func (u *RBACUseCaseImplementation) CheckInsertPermission(ctx context.Context, username, database, schema, table string) (bool, error) {
	return false, errors.New("not implemented")
}
