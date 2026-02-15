package rbac

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *RBACUseCaseImplementation) GetTablePermissions(ctx context.Context, username, database, schema, table string) (*domain.PermissionSet, error) {
	return nil, errors.New("not implemented")
}
