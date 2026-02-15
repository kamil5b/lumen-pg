package dataview

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) GetTableRowCount(ctx context.Context, username, database, schema, table string) (int64, error) {
	// Check if user has SELECT permission
	hasPermission, err := u.rbacRepo.HasSelectPermission(ctx, username, database, schema, table)
	if err != nil {
		return 0, err
	}
	if !hasPermission {
		return 0, domain.ValidationError{
			Field:   "table",
			Message: "user does not have SELECT permission on this table",
		}
	}

	// Get row count from database
	count, err := u.databaseRepo.GetRowCount(ctx, database, schema, table, "")
	if err != nil {
		return 0, err
	}

	return count, nil
}
