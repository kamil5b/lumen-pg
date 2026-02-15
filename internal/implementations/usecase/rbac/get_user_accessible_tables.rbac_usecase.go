package rbac

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *RBACUseCaseImplementation) GetUserAccessibleTables(ctx context.Context, username, database, schema string) ([]domain.AccessibleTable, error) {
	// Get accessible tables for the user
	tables, err := u.rbacRepo.GetAccessibleTables(ctx, username, database, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to get accessible tables: %w", err)
	}

	if tables == nil {
		return []domain.AccessibleTable{}, nil
	}

	return tables, nil
}
