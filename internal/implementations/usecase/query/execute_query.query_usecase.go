package query

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *QueryUseCaseImplementation) ExecuteQuery(ctx context.Context, username, query string, offset, limit int) (*domain.QueryResult, error) {
	// Validate query is not empty
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	// Check if it's a SELECT query
	isSelect, err := u.IsSelectQuery(ctx, query)
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

	// Execute the query with the database repository
	result, err := u.databaseRepo.ExecuteQuery(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("unexpected nil result from database")
	}

	// Apply offset and limit to the results if needed
	if offset > 0 || limit > 0 {
		if offset < 0 {
			offset = 0
		}
		if limit <= 0 {
			limit = 10
		}
		if limit > 1000 {
			limit = 1000 // Hard cap
		}

		// Apply pagination to rows
		if offset < len(result.Rows) {
			endIdx := offset + limit
			if endIdx > len(result.Rows) {
				endIdx = len(result.Rows)
			}
			result.Rows = result.Rows[offset:endIdx]
		} else {
			result.Rows = []map[string]interface{}{}
		}
	}

	return result, nil
}
