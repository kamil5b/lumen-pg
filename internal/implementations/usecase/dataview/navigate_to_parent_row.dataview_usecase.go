package dataview

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) NavigateToParentRow(ctx context.Context, username, database, schema, table, columnName string, value interface{}) (*domain.QueryResult, error) {
	// Check if user has SELECT permission
	hasPermission, err := u.rbacRepo.HasSelectPermission(ctx, username, database, schema, table)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, domain.ValidationError{Message: "no select permission on table"}
	}

	// Build WHERE clause to filter for the parent row
	var whereClause string
	switch v := value.(type) {
	case string:
		whereClause = fmt.Sprintf("%s = '%s'", columnName, v)
	case int, int64, float64:
		whereClause = fmt.Sprintf("%s = %v", columnName, v)
	default:
		whereClause = fmt.Sprintf("%s = '%v'", columnName, v)
	}

	// Build the table data params
	params := domain.TableDataParams{
		Database:    database,
		Schema:      schema,
		Table:       table,
		WhereClause: whereClause,
		Limit:       1,
	}

	// Get parent row from database
	result, err := u.databaseRepo.GetTableData(ctx, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}
