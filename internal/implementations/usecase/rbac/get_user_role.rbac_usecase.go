package rbac

import (
	"context"
	"errors"
)

func (u *RBACUseCaseImplementation) GetUserRole(ctx context.Context, username string) (string, error) {
	return "", errors.New("not implemented")
}
