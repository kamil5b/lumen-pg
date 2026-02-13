package usecase

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// DataViewUseCase defines operations for viewing and interacting with table data
type DataViewUseCase interface {
	// LoadTableData loads data from a table with optional filtering and pagination
	LoadTableData(ctx context.Context, username string, params domain.TableDataParams) (*domain.QueryResult, error)

	// GetTableDataWithCursorPagination loads table data with cursor-based pagination
	GetTableDataWithCursorPagination(ctx context.Context, username, database, schema, table, cursor string, limit int) (*domain.QueryResult, error)

	// FilterTableData filters table data with a WHERE clause
	FilterTableData(ctx context.Context, username, database, schema, table, whereClause string, offset, limit int) (*domain.QueryResult, error)

	// ValidateWhereClause validates a WHERE clause fragment for SQL injection
	ValidateWhereClause(ctx context.Context, whereClause string) (bool, error)

	// SortTableData sorts table data by a column
	SortTableData(ctx context.Context, username, database, schema, table, orderBy, orderDir string, offset, limit int) (*domain.QueryResult, error)

	// GetTableRowCount returns the total count of rows in a table
	GetTableRowCount(ctx context.Context, username, database, schema, table string) (int64, error)

	// GetTableRowCountWithFilter returns the count of rows matching a WHERE clause
	GetTableRowCountWithFilter(ctx context.Context, username, database, schema, table, whereClause string) (int64, error)

	// GetForeignKeyInfo returns information about foreign keys in a table
	GetForeignKeyInfo(ctx context.Context, username, database, schema, table string) ([]domain.ForeignKeyInfo, error)

	// GetPrimaryKeyInfo returns information about primary keys in a table
	GetPrimaryKeyInfo(ctx context.Context, username, database, schema, table string) ([]string, error)

	// GetChildTableReferences returns all child tables that reference a row via foreign key
	GetChildTableReferences(ctx context.Context, username, database, schema, table string, pkValues map[string]interface{}) ([]domain.ChildTableReference, error)

	// GetChildTableRowCount returns the count of rows in a child table that reference a parent row
	GetChildTableRowCount(ctx context.Context, username, database, schema, childTable, parentTable string, fkColumn, pkValue string) (int64, error)

	// NavigateToParentRow navigates to the parent row referenced by a foreign key
	NavigateToParentRow(ctx context.Context, username, database, schema, table, columnName string, value interface{}) (*domain.QueryResult, error)

	// NavigateToChildRows navigates to child rows that reference a parent row
	NavigateToChildRows(ctx context.Context, username, database, schema, childTable, parentTable string, fkColumn, pkColumn, pkValue string) (*domain.QueryResult, error)

	// IsTableReadOnly checks if a table is read-only for the user
	IsTableReadOnly(ctx context.Context, username, database, schema, table string) (bool, error)
}
