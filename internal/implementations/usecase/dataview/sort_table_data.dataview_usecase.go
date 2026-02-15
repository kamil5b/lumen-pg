package dataview

import (
	"context"
	"strings"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) SortTableData(ctx context.Context, username, database, schema, table, orderBy, orderDir string, offset, limit int) (*domain.QueryResult, error) {
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

	// Validate order direction
	orderDir = strings.ToUpper(strings.TrimSpace(orderDir))
	if orderDir != "ASC" && orderDir != "DESC" {
		return nil, domain.ValidationError{
			Field:   "orderDir",
			Message: "order direction must be ASC or DESC",
		}
	}

	// Build the table data params
	params := domain.TableDataParams{
		Database: database,
		Schema:   schema,
		Table:    table,
		OrderBy:  orderBy,
		OrderDir: orderDir,
		Offset:   offset,
		Limit:    limit,
	}

	// Get sorted table data from database
	result, err := u.databaseRepo.GetTableData(ctx, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}
