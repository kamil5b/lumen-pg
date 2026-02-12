package interfaces

import (
	"context"
	"database/sql"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// ConnectionRepository handles database connection management
type ConnectionRepository interface {
	// ValidateConnectionString validates the format of a connection string
	ValidateConnectionString(connStr string) error

	// TestConnection tests connectivity to a PostgreSQL instance
	TestConnection(ctx context.Context, connStr string) error

	// GetConnection returns a database connection for a user
	GetConnection(ctx context.Context, username, password, database string) (*sql.DB, error)

	// ProbeFirstAccessible probes connection to the first accessible database and table for a user
	ProbeFirstAccessible(ctx context.Context, username, password string, metadata *domain.GlobalMetadata) (string, error)
}

// MetadataRepository handles database metadata operations
type MetadataRepository interface {
	// LoadGlobalMetadata loads all metadata about databases, schemas, tables, and roles
	LoadGlobalMetadata(ctx context.Context, db *sql.DB) (*domain.GlobalMetadata, error)

	// LoadDatabaseMetadata loads metadata for a specific database
	LoadDatabaseMetadata(ctx context.Context, db *sql.DB, database string) (*domain.DatabaseMetadata, error)

	// LoadRolePermissions loads permissions for all roles
	LoadRolePermissions(ctx context.Context, db *sql.DB) (map[string]*domain.RolePermissions, error)

	// GetTableMetadata gets metadata for a specific table
	GetTableMetadata(ctx context.Context, db *sql.DB, schema, table string) (*domain.TableMetadata, error)

	// GetERDData generates ERD data for a schema
	GetERDData(ctx context.Context, db *sql.DB, database, schema string) (*domain.ERDData, error)
}

// QueryRepository handles query execution
type QueryRepository interface {
	// ExecuteQuery executes a SQL query
	ExecuteQuery(ctx context.Context, db *sql.DB, req domain.QueryRequest) (*domain.QueryResult, error)

	// ExecuteMultipleQueries executes multiple SQL queries separated by semicolons
	ExecuteMultipleQueries(ctx context.Context, db *sql.DB, sql string) ([]*domain.QueryResult, error)

	// GetTableData retrieves paginated table data
	GetTableData(ctx context.Context, db *sql.DB, req domain.TableDataRequest) (*domain.TableDataResult, error)

	// GetReferencingTables gets tables that reference a primary key
	GetReferencingTables(ctx context.Context, db *sql.DB, schema, table, pkColumn string, pkValue interface{}) (map[string]int64, error)
}

// TransactionRepository handles transaction operations
type TransactionRepository interface {
	// BeginTransaction starts a new database transaction
	BeginTransaction(ctx context.Context, db *sql.DB) (*sql.Tx, error)

	// CommitTransaction commits a transaction
	CommitTransaction(ctx context.Context, tx *sql.Tx) error

	// RollbackTransaction rolls back a transaction
	RollbackTransaction(ctx context.Context, tx *sql.Tx) error

	// ExecuteOperations executes buffered operations within a transaction
	ExecuteOperations(ctx context.Context, tx *sql.Tx, operations []domain.Operation) error
}
