package dataview

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) GetChildTableRowCount(ctx context.Context, username, database, schema, childTable, parentTable string, fkColumn, pkValue string) (int64, error) {
	// Check if user has SELECT permission on child table
	hasPermission, err := u.rbacRepo.HasSelectPermission(ctx, username, database, schema, childTable)
	if err != nil {
		return 0, err
	}
	if !hasPermission {
		return 0, domain.ValidationError{Message: "no select permission on child table"}
	}

	// Check if user has SELECT permission on parent table
	hasPermissionParent, err := u.rbacRepo.HasSelectPermission(ctx, username, database, schema, parentTable)
	if err != nil {
		return 0, err
	}
	if !hasPermissionParent {
		return 0, domain.ValidationError{Message: "no select permission on parent table"}
	}

	// Build WHERE clause to filter by foreign key
	whereClause := fmt.Sprintf("%s = '%s'", fkColumn, pkValue)

	// Get row count with filter from database
	return u.databaseRepo.GetRowCount(ctx, database, schema, childTable, whereClause)
}
