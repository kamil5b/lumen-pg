package rbac

import (
	"context"
	"errors"
)

func (u *RBACUseCaseImplementation) GetUserAccessibleDatabases(ctx context.Context, username string) ([]string, error) {
	return nil, errors.New("not implemented")
}
