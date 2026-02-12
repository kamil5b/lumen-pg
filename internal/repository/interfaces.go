package repository

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// ConnectionRepository handles database connection operations.
type ConnectionRepository interface {
	// ValidateConnectionString validates the format of a connection string.
	ValidateConnectionString(connStr string) error

	// ParseConnectionString parses a connection string into its components.
	ParseConnectionString(connStr string) (*domain.ConnectionConfig, error)

	// TestConnection tests connectivity to a PostgreSQL instance.
	TestConnection(ctx context.Context, connStr string) error

	// ProbeConnection probes connection to the first accessible database, schema, and table for a user.
	ProbeConnection(ctx context.Context, username, password, host string, port int, sslMode string, accessibleDBs []string) (*domain.Session, error)

	// Connect creates a connection to a specific database.
	Connect(ctx context.Context, config *domain.ConnectionConfig) (interface{}, error)
}

// MetadataRepository handles database metadata operations.
type MetadataRepository interface {
	// LoadDatabases fetches all databases from the PostgreSQL instance.
	LoadDatabases(ctx context.Context) ([]domain.Database, error)

	// LoadSchemas fetches all schemas for a database.
	LoadSchemas(ctx context.Context, database string) ([]domain.Schema, error)

	// LoadTables fetches all tables for a schema.
	LoadTables(ctx context.Context, database, schema string) ([]domain.Table, error)

	// LoadColumns fetches all columns for a table.
	LoadColumns(ctx context.Context, database, schema, table string) ([]domain.Column, error)

	// LoadForeignKeys fetches all foreign key relationships for a table.
	LoadForeignKeys(ctx context.Context, database, schema, table string) ([]domain.ForeignKey, error)

	// LoadRoles fetches all PostgreSQL roles and their permissions.
	LoadRoles(ctx context.Context) ([]domain.Role, error)

	// LoadAllMetadata fetches all metadata including databases, schemas, tables, columns, relations, and roles.
	LoadAllMetadata(ctx context.Context) (*domain.Metadata, error)

	// GetAccessibleResources returns accessible resources for a specific role.
	GetAccessibleResources(ctx context.Context, roleName string) (*domain.Role, error)

	// GenerateERDData generates ERD data for a specific schema.
	GenerateERDData(ctx context.Context, database, schema string) (*domain.ERDData, error)
}

// QueryRepository handles SQL query execution operations.
type QueryRepository interface {
	// ExecuteQuery executes a single SQL query and returns the result.
	ExecuteQuery(ctx context.Context, database, query string, params ...interface{}) (*domain.QueryResult, error)

	// ExecuteQueries executes multiple SQL queries in sequence.
	ExecuteQueries(ctx context.Context, database string, queries []string) ([]domain.QueryResult, error)

	// LoadTableData loads table data with cursor pagination.
	LoadTableData(ctx context.Context, database, schema, table string, cursor *domain.Cursor, orderBy string, orderDir string, whereClause string) (*domain.CursorPage, error)

	// GetTotalRowCount gets the total row count for a table with optional WHERE clause.
	GetTotalRowCount(ctx context.Context, database, schema, table string, whereClause string) (int64, error)

	// GetReferencingTables returns tables that reference a primary key value.
	GetReferencingTables(ctx context.Context, database, schema, table, pkColumn string, pkValue interface{}) ([]domain.ReferencingTable, error)

	// ExecuteTransaction executes buffered operations atomically.
	ExecuteTransaction(ctx context.Context, database string, operations []domain.BufferedOperation) error
}
