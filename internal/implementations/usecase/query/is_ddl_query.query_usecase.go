package query

import (
	"context"
	"strings"
)

func (u *QueryUseCaseImplementation) IsDDLQuery(ctx context.Context, query string) (bool, error) {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return false, nil
	}

	// Get the first word and convert to uppercase
	firstWord := strings.ToUpper(strings.Fields(trimmed)[0])

	// Check if it's a DDL statement
	ddlKeywords := map[string]bool{
		"CREATE":   true,
		"ALTER":    true,
		"DROP":     true,
		"TRUNCATE": true,
	}

	return ddlKeywords[firstWord], nil
}
