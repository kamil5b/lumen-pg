package usecase

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// QueryUseCase defines operations for executing SQL queries
type QueryUseCase interface {
	// ExecuteQuery executes a single SQL query
	ExecuteQuery(ctx context.Context, username, query string, offset, limit int) (*domain.QueryResult, error)

	// ExecuteMultipleQueries executes multiple SQL queries separated by semicolons
	ExecuteMultipleQueries(ctx context.Context, username, queries string) ([]domain.QueryResult, error)

	// ExecuteQueryWithPagination executes a query with offset pagination
	ExecuteQueryWithPagination(ctx context.Context, username string, params domain.QueryParams) (*domain.QueryResult, error)

	// SplitQueries splits a multi-query string by semicolons
	SplitQueries(ctx context.Context, queries string) ([]string, error)

	// ValidateQuery validates SQL query syntax
	ValidateQuery(ctx context.Context, query string) (bool, error)

	// IsSelectQuery checks if a query is a SELECT statement
	IsSelectQuery(ctx context.Context, query string) (bool, error)

	// IsDDLQuery checks if a query is a DDL statement (CREATE, ALTER, DROP)
	IsDDLQuery(ctx context.Context, query string) (bool, error)

	// IsDMLQuery checks if a query is a DML statement (INSERT, UPDATE, DELETE)
	IsDMLQuery(ctx context.Context, query string) (bool, error)

	// GetQueryAffectedRowCount returns the number of rows affected by a query
	GetQueryAffectedRowCount(ctx context.Context, result *domain.QueryResult) int64
}
