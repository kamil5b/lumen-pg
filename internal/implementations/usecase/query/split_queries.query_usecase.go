package query

import (
	"context"
	"strings"
)

func (u *QueryUseCaseImplementation) SplitQueries(ctx context.Context, queries string) ([]string, error) {
	if queries == "" {
		return []string{}, nil
	}

	// Split by semicolon
	parts := strings.Split(queries, ";")

	var result []string
	for _, part := range parts {
		// Trim whitespace and skip empty strings
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result, nil
}
