package query

import (
	"context"
	"strings"
)

func (u *QueryUseCaseImplementation) IsSelectQuery(ctx context.Context, query string) (bool, error) {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return false, nil
	}

	// Get the first word and convert to uppercase
	firstWord := strings.Fields(trimmed)[0]
	return strings.EqualFold(firstWord, "SELECT"), nil
}
