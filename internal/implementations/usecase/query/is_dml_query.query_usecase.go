package query

import (
	"context"
	"strings"
)

func (u *QueryUseCaseImplementation) IsDMLQuery(ctx context.Context, query string) (bool, error) {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return false, nil
	}

	// Get the first word and convert to uppercase
	firstWord := strings.ToUpper(strings.Fields(trimmed)[0])

	// Check if it's a DML statement
	dmlKeywords := map[string]bool{
		"INSERT": true,
		"UPDATE": true,
		"DELETE": true,
	}

	return dmlKeywords[firstWord], nil
}
