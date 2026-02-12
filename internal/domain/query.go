package domain

import (
	"errors"
	"strings"
)

// SplitQueries splits a SQL string into individual queries by semicolons,
// respecting quoted strings.
func SplitQueries(sql string) []string {
	var queries []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false

	for i := 0; i < len(sql); i++ {
		ch := sql[i]

		switch {
		case ch == '\'' && !inDoubleQuote:
			inSingleQuote = !inSingleQuote
			current.WriteByte(ch)
		case ch == '"' && !inSingleQuote:
			inDoubleQuote = !inDoubleQuote
			current.WriteByte(ch)
		case ch == ';' && !inSingleQuote && !inDoubleQuote:
			q := strings.TrimSpace(current.String())
			if q != "" {
				queries = append(queries, q)
			}
			current.Reset()
		default:
			current.WriteByte(ch)
		}
	}

	q := strings.TrimSpace(current.String())
	if q != "" {
		queries = append(queries, q)
	}

	return queries
}

// ValidateWhereClause validates a WHERE clause fragment for safety.
func ValidateWhereClause(clause string) error {
	if clause == "" {
		return nil
	}

	lower := strings.ToLower(strings.TrimSpace(clause))

	// Check for dangerous SQL keywords that shouldn't appear in a WHERE clause
	dangerousKeywords := []string{
		"drop ", "alter ", "create ", "truncate ",
		"grant ", "revoke ", "insert ", "update ",
		"delete ", "--", "/*",
	}

	for _, kw := range dangerousKeywords {
		if strings.Contains(lower, kw) {
			return ErrSQLInjectionDetected
		}
	}

	return nil
}

var ErrSQLInjectionDetected = errors.New("potential SQL injection detected")
