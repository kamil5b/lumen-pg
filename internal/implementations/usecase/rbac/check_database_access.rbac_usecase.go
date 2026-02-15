package rbac

import (
	"context"
	"errors"
)

func (u *RBACUseCaseImplementation) CheckDatabaseAccess(ctx context.Context, username, database string) (bool, error) {
	return false, errors.New("not implemented")
}
