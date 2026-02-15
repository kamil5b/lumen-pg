package query

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *QueryUseCaseImplementation) ExecuteQueryWithPagination(ctx context.Context, username string, params domain.QueryParams) (*domain.QueryResult, error) {
	// Validate query is not empty
	if params.Query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	// Check if it's a SELECT query
	isSelect, err := u.IsSelectQuery(ctx, params.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to check query type: %w", err)
	}

	if !isSelect {
		return nil, fmt.Errorf("only SELECT queries are allowed")
	}

	// Check RBAC permissions for SELECT
	hasPermission, err := u.rbacRepo.HasSelectPermission(ctx, username, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to check permissions: %w", err)
	}

	if !hasPermission {
		return nil, fmt.Errorf("access denied: user does not have SELECT permission")
	}

	// Apply hard limit cap
	limit := params.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 1000 {
		limit = 1000 // Hard cap at 1000 rows
	}

	// Ensure offset is not negative
	offset := params.Offset
	if offset < 0 {
		offset = 0
	}

	// Update params with corrected values
	params.Offset = offset
	params.Limit = limit

	// Execute the query with pagination
	result, err := u.databaseRepo.ExecuteQueryWithPagination(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("unexpected nil result from database")
	}

	return result, nil
}
