package dataview

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) LoadTableData(ctx context.Context, username string, params domain.TableDataParams) (*domain.QueryResult, error) {
	// Check if user has SELECT permission on the table
	hasPermission, err := u.rbacRepo.HasSelectPermission(ctx, username, params.Database, params.Schema, params.Table)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, domain.ValidationError{
			Field:   "table",
			Message: "user does not have SELECT permission on this table",
		}
	}

	// Get table data from database
	result, err := u.databaseRepo.GetTableData(ctx, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}
