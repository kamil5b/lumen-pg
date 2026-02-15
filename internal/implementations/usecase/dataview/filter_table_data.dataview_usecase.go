package dataview

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) FilterTableData(ctx context.Context, username, database, schema, table, whereClause string, offset, limit int) (*domain.QueryResult, error) {
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

	// Validate the WHERE clause for SQL injection
	valid, err := u.ValidateWhereClause(ctx, whereClause)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, domain.ValidationError{
			Field:   "whereClause",
			Message: "WHERE clause contains invalid or malicious patterns",
		}
	}

	// Build the table data params
	params := domain.TableDataParams{
		Database:    database,
		Schema:      schema,
		Table:       table,
		WhereClause: whereClause,
		Offset:      offset,
		Limit:       limit,
	}

	// Get filtered table data from database
	result, err := u.databaseRepo.GetTableData(ctx, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}
