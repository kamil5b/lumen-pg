package usecase

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// AuthUsecase handles authentication operations.
type AuthUsecase interface {
	// Login validates credentials, probes connection, and creates a session.
	Login(ctx context.Context, username, password string) (*domain.Session, error)

	// ValidateSession validates an existing session.
	ValidateSession(session *domain.Session) error

	// Logout clears a user session.
	Logout(ctx context.Context, username string) error

	// ReAuthenticate re-authenticates using encrypted password from cookie.
	ReAuthenticate(ctx context.Context, username, encryptedPassword string) error

	// GetAccessibleResources returns the accessible resources for a user.
	GetAccessibleResources(ctx context.Context, username string) (*domain.Role, error)
}

// MetadataUsecase handles metadata operations.
type MetadataUsecase interface {
	// InitializeMetadata loads and caches all metadata from the PostgreSQL instance.
	InitializeMetadata(ctx context.Context, connStr string) error

	// GetMetadata returns the cached metadata.
	GetMetadata() *domain.Metadata

	// RefreshMetadata reloads metadata from the database.
	RefreshMetadata(ctx context.Context) error

	// GetRBACMapping returns the role-based access control mapping.
	GetRBACMapping() map[string]*domain.Role

	// GetERDData generates ERD data for a specific schema.
	GetERDData(ctx context.Context, database, schema string) (*domain.ERDData, error)
}

// DataExplorerUsecase handles main view data operations.
type DataExplorerUsecase interface {
	// LoadTableData loads table data with cursor pagination.
	LoadTableData(ctx context.Context, database, schema, table string, cursor *domain.Cursor, orderBy string, orderDir string, whereClause string) (*domain.CursorPage, error)

	// GetTotalRowCount returns the total row count for a table.
	GetTotalRowCount(ctx context.Context, database, schema, table string, whereClause string) (int64, error)

	// GetReferencingTables returns tables that reference a primary key value.
	GetReferencingTables(ctx context.Context, database, schema, table, pkColumn string, pkValue interface{}) ([]domain.ReferencingTable, error)

	// NavigateToForeignKey navigates to the parent table of a foreign key.
	NavigateToForeignKey(ctx context.Context, database, schema, table, column string, value interface{}) (*domain.CursorPage, error)
}

// QueryUsecase handles manual query execution.
type QueryUsecase interface {
	// ExecuteQuery executes a single SQL query.
	ExecuteQuery(ctx context.Context, database, query string) (*domain.QueryResult, error)

	// ExecuteQueries executes multiple SQL queries.
	ExecuteQueries(ctx context.Context, database, sql string) ([]domain.QueryResult, error)

	// ExecuteQueryWithPagination executes a SELECT query with offset pagination.
	ExecuteQueryWithPagination(ctx context.Context, database, query string, offset, limit int) (*domain.QueryResult, error)
}

// TransactionUsecase handles transaction operations.
type TransactionUsecase interface {
	// StartTransaction starts a new transaction for a user session.
	StartTransaction(ctx context.Context, sessionID string) error

	// AddOperation adds a buffered operation to the active transaction.
	AddOperation(ctx context.Context, sessionID string, op domain.BufferedOperation) error

	// Commit commits the active transaction.
	Commit(ctx context.Context, sessionID string) error

	// Rollback rolls back the active transaction.
	Rollback(ctx context.Context, sessionID string) error

	// GetTransaction returns the active transaction for a session.
	GetTransaction(sessionID string) (*domain.Transaction, error)
}
