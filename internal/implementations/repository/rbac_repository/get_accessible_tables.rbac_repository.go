package rbac_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (r *RBACRepositoryImplementation) GetAccessibleTables(ctx context.Context, role, database, schema string) ([]domain.AccessibleTable, error) {
	return nil, errors.New("not implemented")
}
