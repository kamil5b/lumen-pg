package interfaces

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// AuthService handles authentication and session management
type AuthService interface {
	// Login authenticates a user and creates a session
	Login(ctx context.Context, input domain.LoginInput) (*domain.Session, error)

	// ValidateSession validates a session token
	ValidateSession(ctx context.Context, sessionToken string) (*domain.Session, error)

	// Logout invalidates a session
	Logout(ctx context.Context, sessionToken string) error

	// EncryptPassword encrypts a password for cookie storage
	EncryptPassword(password string) (string, error)

	// DecryptPassword decrypts a password from cookie storage
	DecryptPassword(encrypted string) (string, error)
}

// MetadataService handles metadata operations
type MetadataService interface {
	// InitializeMetadata initializes global metadata cache
	InitializeMetadata(ctx context.Context, superAdminConnStr string) (*domain.GlobalMetadata, error)

	// RefreshMetadata refreshes the global metadata cache
	RefreshMetadata(ctx context.Context) (*domain.GlobalMetadata, error)

	// GetAccessibleResources gets resources accessible to a role
	GetAccessibleResources(ctx context.Context, roleName string) (*domain.RolePermissions, error)

	// GetERDData gets ERD data for a schema
	GetERDData(ctx context.Context, username, database, schema string) (*domain.ERDData, error)
}

// QueryService handles query execution
type QueryService interface {
	// ExecuteQuery executes a SQL query
	ExecuteQuery(ctx context.Context, username, database string, req domain.QueryRequest) (*domain.QueryResult, error)

	// ExecuteMultipleQueries executes multiple SQL queries
	ExecuteMultipleQueries(ctx context.Context, username, database, sql string) ([]*domain.QueryResult, error)

	// ValidateWhereClause validates a WHERE clause for safety
	ValidateWhereClause(whereClause string) error

	// SplitQueries splits SQL by semicolons intelligently
	SplitQueries(sql string) ([]string, error)
}

// DataExplorerService handles main view data operations
type DataExplorerService interface {
	// GetTableData gets paginated table data
	GetTableData(ctx context.Context, username, database string, req domain.TableDataRequest) (*domain.TableDataResult, error)

	// GetReferencingTables gets tables referencing a primary key
	GetReferencingTables(ctx context.Context, username, database, schema, table, pkColumn string, pkValue interface{}) (map[string]int64, error)

	// NavigateToForeignKey navigates to a parent table via foreign key
	NavigateToForeignKey(ctx context.Context, username, database, schema, table, fkColumn string, fkValue interface{}) (*domain.TableDataResult, error)
}

// TransactionService handles transaction operations
type TransactionService interface {
	// StartTransaction starts a new transaction for a user
	StartTransaction(ctx context.Context, username, database string) (*domain.TransactionState, error)

	// BufferOperation adds an operation to the transaction buffer
	BufferOperation(ctx context.Context, username string, op domain.Operation) error

	// CommitTransaction commits the active transaction
	CommitTransaction(ctx context.Context, username string) error

	// RollbackTransaction rolls back the active transaction
	RollbackTransaction(ctx context.Context, username string) error

	// GetTransactionState gets the current transaction state for a user
	GetTransactionState(ctx context.Context, username string) (*domain.TransactionState, error)

	// CheckTransactionTimeout checks if a transaction has expired
	CheckTransactionTimeout(ctx context.Context, username string) (bool, error)
}
