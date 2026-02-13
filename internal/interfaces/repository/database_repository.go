package repository

import (
	"context"
	"database/sql"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// DatabaseRepository defines operations for database connections
type DatabaseRepository interface {
	// Connect establishes a connection to a PostgreSQL database
	Connect(ctx context.Context, connString string) error

	// Disconnect closes the database connection
	Disconnect(ctx context.Context) error

	// TestConnection verifies connectivity to a database
	TestConnection(ctx context.Context, connString string) error

	// GetConnection returns the active database connection
	GetConnection() *sql.DB

	// ExecuteQuery executes a SQL query and returns results
	ExecuteQuery(ctx context.Context, query string, args ...interface{}) (*domain.QueryResult, error)

	// ExecuteQueryWithPagination executes a query with pagination
	ExecuteQueryWithPagination(ctx context.Context, params domain.QueryParams) (*domain.QueryResult, error)

	// ExecuteMultipleQueries executes multiple SQL queries separated by semicolons
	ExecuteMultipleQueries(ctx context.Context, queries string) ([]domain.QueryResult, error)

	// BeginTransaction starts a new transaction
	BeginTransaction(ctx context.Context) (*sql.Tx, error)

	// CommitTransaction commits an active transaction
	CommitTransaction(ctx context.Context, tx *sql.Tx) error

	// RollbackTransaction rolls back an active transaction
	RollbackTransaction(ctx context.Context, tx *sql.Tx) error

	// GetDatabases retrieves list of available databases
	GetDatabases(ctx context.Context) ([]string, error)

	// GetSchemas retrieves list of schemas in a database
	GetSchemas(ctx context.Context, database string) ([]string, error)

	// GetTables retrieves list of tables in a schema
	GetTables(ctx context.Context, database, schema string) ([]string, error)

	// GetTableMetadata retrieves detailed metadata for a table
	GetTableMetadata(ctx context.Context, database, schema, table string) (*domain.TableMetadata, error)

	// GetDatabaseMetadata retrieves complete metadata for a database
	GetDatabaseMetadata(ctx context.Context, database string) (*domain.DatabaseMetadata, error)

	// GetTableData retrieves data from a table with optional filtering and pagination
	GetTableData(ctx context.Context, params domain.TableDataParams) (*domain.QueryResult, error)

	// InsertRow inserts a new row into a table
	InsertRow(ctx context.Context, database, schema, table string, values map[string]interface{}) error

	// UpdateRow updates a row in a table
	UpdateRow(ctx context.Context, database, schema, table string, pkValues map[string]interface{}, values map[string]interface{}) error

	// DeleteRow deletes a row from a table
	DeleteRow(ctx context.Context, database, schema, table string, pkValues map[string]interface{}) error

	// GetRowCount retrieves the total count of rows in a table (with optional WHERE clause)
	GetRowCount(ctx context.Context, database, schema, table, whereClause string) (int64, error)
}
