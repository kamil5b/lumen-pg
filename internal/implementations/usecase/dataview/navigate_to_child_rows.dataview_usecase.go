package dataview

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) NavigateToChildRows(ctx context.Context, username, database, schema, childTable, parentTable string, fkColumn, pkColumn, pkValue string) (*domain.QueryResult, error) {
	// Check if user has SELECT permission on child table
	hasPermission, err := u.rbacRepo.HasSelectPermission(ctx, username, database, schema, childTable)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, domain.ValidationError{Message: "no select permission on child table"}
	}

	// Build WHERE clause to filter child rows by foreign key value
	whereClause := fmt.Sprintf("%s = '%s'", fkColumn, pkValue)

	// Build the table data params
	params := domain.TableDataParams{
		Database:    database,
		Schema:      schema,
		Table:       childTable,
		WhereClause: whereClause,
	}

	// Get child rows from database
	result, err := u.databaseRepo.GetTableData(ctx, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}
