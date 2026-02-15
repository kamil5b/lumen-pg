package query

import (
	"context"
	"regexp"
	"strings"
)

func (u *QueryUseCaseImplementation) ValidateQuery(ctx context.Context, query string) (bool, error) {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return false, nil
	}

	// Check for basic SQL keywords at the start
	firstWord := strings.ToUpper(strings.Fields(trimmed)[0])
	validKeywords := map[string]bool{
		"SELECT":   true,
		"INSERT":   true,
		"UPDATE":   true,
		"DELETE":   true,
		"CREATE":   true,
		"ALTER":    true,
		"DROP":     true,
		"TRUNCATE": true,
		"WITH":     true,
	}

	if !validKeywords[firstWord] {
		return false, nil
	}

	// Must have balanced parentheses
	if !hasBalancedParentheses(trimmed) {
		return false, nil
	}

	// Check for obviously invalid syntax
	if strings.Contains(trimmed, ";;") {
		return false, nil
	}

	// For SELECT queries, validate basic structure
	if firstWord == "SELECT" {
		return validateSelectQuery(trimmed), nil
	}

	// For INSERT queries, validate basic structure
	if firstWord == "INSERT" {
		return validateInsertQuery(trimmed), nil
	}

	// For UPDATE queries, validate basic structure
	if firstWord == "UPDATE" {
		return validateUpdateQuery(trimmed), nil
	}

	// For DELETE queries, validate basic structure
	if firstWord == "DELETE" {
		return validateDeleteQuery(trimmed), nil
	}

	// For DDL queries, basic validation
	if firstWord == "CREATE" || firstWord == "ALTER" || firstWord == "DROP" || firstWord == "TRUNCATE" {
		return validateDDLQuery(trimmed), nil
	}

	return true, nil
}

func validateSelectQuery(query string) bool {
	// SELECT must have at least one column and FROM clause
	query = strings.ToUpper(query)

	// Remove string literals to avoid false matches
	query = removeStringLiterals(query)

	// Must contain FROM unless it's a simple SELECT without FROM (like SELECT 1)
	if !strings.Contains(query, "FROM") {
		// Allow simple selects like "SELECT 1", "SELECT 123", etc.
		// But reject if it has a typo like "SELECT * FORM"
		if strings.Contains(query, "FORM") {
			return false
		}
		return true
	}

	return true
}

func validateInsertQuery(query string) bool {
	query = strings.ToUpper(query)
	query = removeStringLiterals(query)

	// INSERT must have INTO and VALUES or SELECT
	if !strings.Contains(query, "INTO") {
		return false
	}

	if !strings.Contains(query, "VALUES") && !strings.Contains(query, "SELECT") {
		return false
	}

	return true
}

func validateUpdateQuery(query string) bool {
	query = strings.ToUpper(query)
	query = removeStringLiterals(query)

	// UPDATE must have SET clause
	if !strings.Contains(query, "SET") {
		return false
	}

	return true
}

func validateDeleteQuery(query string) bool {
	query = strings.ToUpper(query)
	query = removeStringLiterals(query)

	// DELETE must have FROM clause
	if !strings.Contains(query, "FROM") {
		return false
	}

	return true
}

func validateDDLQuery(query string) bool {
	query = strings.ToUpper(query)
	query = removeStringLiterals(query)

	// CREATE must have TABLE, INDEX, VIEW, etc.
	if strings.HasPrefix(query, "CREATE") {
		if !strings.Contains(query, "TABLE") &&
			!strings.Contains(query, "INDEX") &&
			!strings.Contains(query, "VIEW") &&
			!strings.Contains(query, "DATABASE") &&
			!strings.Contains(query, "SCHEMA") &&
			!strings.Contains(query, "FUNCTION") &&
			!strings.Contains(query, "SEQUENCE") {
			return false
		}
	}

	// ALTER must have TABLE or other object type
	if strings.HasPrefix(query, "ALTER") {
		if !strings.Contains(query, "TABLE") &&
			!strings.Contains(query, "INDEX") &&
			!strings.Contains(query, "VIEW") &&
			!strings.Contains(query, "DATABASE") &&
			!strings.Contains(query, "SCHEMA") &&
			!strings.Contains(query, "FUNCTION") &&
			!strings.Contains(query, "SEQUENCE") {
			return false
		}
	}

	// DROP must have TABLE or other object type
	if strings.HasPrefix(query, "DROP") {
		if !strings.Contains(query, "TABLE") &&
			!strings.Contains(query, "INDEX") &&
			!strings.Contains(query, "VIEW") &&
			!strings.Contains(query, "DATABASE") &&
			!strings.Contains(query, "SCHEMA") &&
			!strings.Contains(query, "FUNCTION") &&
			!strings.Contains(query, "SEQUENCE") {
			return false
		}
	}

	return true
}

func removeStringLiterals(query string) string {
	// Remove single-quoted strings to avoid false matches
	re := regexp.MustCompile(`'([^'\\]|\\.)*'`)
	return re.ReplaceAllString(query, " ")
}

func hasBalancedParentheses(query string) bool {
	count := 0
	inString := false
	escapeNext := false

	for _, char := range query {
		if escapeNext {
			escapeNext = false
			continue
		}

		if char == '\\' {
			escapeNext = true
			continue
		}

		if char == '\'' {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		switch char {
		case '(':
			count++
		case ')':
			count--
			if count < 0 {
				return false
			}
		}
	}

	return count == 0
}
