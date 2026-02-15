package query

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *QueryUseCaseImplementation) ExecuteMultipleQueries(ctx context.Context, username, queries string) ([]domain.QueryResult, error) {
	// Validate queries is not empty
	if queries == "" {
		return nil, fmt.Errorf("queries cannot be empty")
	}

	// Split queries by semicolon to validate them
	splitQueries, err := u.SplitQueries(ctx, queries)
	if err != nil {
		return nil, fmt.Errorf("failed to split queries: %w", err)
	}

	if len(splitQueries) == 0 {
		return nil, fmt.Errorf("no valid queries found")
	}

	// Check RBAC permissions for SELECT
	hasPermission, err := u.rbacRepo.HasSelectPermission(ctx, username, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to check permissions: %w", err)
	}

	if !hasPermission {
		return nil, fmt.Errorf("access denied: user does not have SELECT permission")
	}

	// Validate all queries are SELECT
	for _, query := range splitQueries {
		isSelect, err := u.IsSelectQuery(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("failed to check query type: %w", err)
		}

		if !isSelect {
			return nil, fmt.Errorf("only SELECT queries are allowed, got: %s", query)
		}
	}

	// Execute multiple queries using database repository
	results, err := u.databaseRepo.ExecuteMultipleQueries(ctx, queries)
	if err != nil {
		return nil, fmt.Errorf("failed to execute queries: %w", err)
	}

	if results == nil {
		return nil, fmt.Errorf("unexpected nil result from database")
	}

	return results, nil
}
