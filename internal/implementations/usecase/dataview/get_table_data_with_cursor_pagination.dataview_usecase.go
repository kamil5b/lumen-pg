package dataview

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) GetTableDataWithCursorPagination(ctx context.Context, username, database, schema, table, cursor string, limit int) (*domain.QueryResult, error) {
	// Check if user has SELECT permission
	hasPermission, err := u.rbacRepo.HasSelectPermission(ctx, username, database, schema, table)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, domain.ValidationError{
			Field:   "table",
			Message: "user does not have SELECT permission on this table",
		}
	}

	// Build the table data params with cursor
	params := domain.TableDataParams{
		Database: database,
		Schema:   schema,
		Table:    table,
		Cursor:   cursor,
		Limit:    limit,
	}

	// Get table data with cursor pagination from database
	result, err := u.databaseRepo.GetTableData(ctx, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}
