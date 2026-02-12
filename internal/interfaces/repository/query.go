package repository

import (
	"context"
	"github.com/kamil5b/lumen-pg/internal/domain"
)

// QueryRepository handles SQL query execution
type QueryRepository interface {
	// ExecuteQuery executes a SQL query with parameters (for SELECT)
	ExecuteQuery(ctx context.Context, query string, params ...interface{}) (*domain.QueryResult, error)
	
	// ExecuteDML executes a DML query (INSERT, UPDATE, DELETE)
	ExecuteDML(ctx context.Context, query string, params ...interface{}) (*domain.QueryResult, error)
	
	// ExecuteDDL executes a DDL query (CREATE, ALTER, DROP)
	ExecuteDDL(ctx context.Context, query string) error
	
	// ExecuteMultiple executes multiple queries separated by semicolons
	ExecuteMultiple(ctx context.Context, queries string) ([]*domain.QueryResult, error)
}
