package rbac

import (
	"context"
	"errors"
)

func (u *RBACUseCaseImplementation) GetUserAccessibleSchemas(ctx context.Context, username, database string) ([]string, error) {
	return nil, errors.New("not implemented")
}
