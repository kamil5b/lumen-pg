package rbac

import (
	"context"
	"errors"
)

func (u *RBACUseCaseImplementation) CheckSchemaAccess(ctx context.Context, username, database, schema string) (bool, error) {
	return false, errors.New("not implemented")
}
