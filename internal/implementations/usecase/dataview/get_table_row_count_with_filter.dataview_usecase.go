package dataview

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) GetTableRowCountWithFilter(ctx context.Context, username, database, schema, table, whereClause string) (int64, error) {
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

	// Validate the WHERE clause for SQL injection
	valid, err := u.ValidateWhereClause(ctx, whereClause)
	if err != nil {
		return 0, err
	}
	if !valid {
		return 0, domain.ValidationError{
			Field:   "whereClause",
			Message: "WHERE clause contains invalid or malicious patterns",
		}
	}

	// Get row count with filter from database
	count, err := u.databaseRepo.GetRowCount(ctx, database, schema, table, whereClause)
	if err != nil {
		return 0, err
	}

	return count, nil
}
