package rbac

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *RBACUseCaseImplementation) GetUserAccessibleTables(ctx context.Context, username, database, schema string) ([]domain.AccessibleTable, error) {
	return nil, errors.New("not implemented")
}
