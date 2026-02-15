package dataview

import (
	"context"
	"regexp"
	"strings"
)

func (u *DataViewUseCaseImplementation) ValidateWhereClause(ctx context.Context, whereClause string) (bool, error) {
	// Check for SQL injection patterns
	dangerousPatterns := []string{
		`(?i)'\s*OR\s*'`,         // ' OR '
		`(?i)'\s*OR\s*1\s*=\s*1`, // ' OR 1=1
		`(?i)--`,                 // SQL comments
		`(?i)/\*.*\*/`,           // Multi-line comments
		`(?i);\s*DROP`,           // DROP statements
		`(?i);\s*DELETE`,         // DELETE statements
		`(?i)UNION\s+SELECT`,     // UNION SELECT
		`(?i)xp_`,                // Extended stored procedures
		`(?i)sp_`,                // System stored procedures
		`(?i)exec\s*\(`,          // EXEC function
		`(?i)execute\s*\(`,       // EXECUTE function
	}

	for _, pattern := range dangerousPatterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(whereClause) {
			return false, nil
		}
	}

	// Additional check: ensure the clause is not empty
	if strings.TrimSpace(whereClause) == "" {
		return false, nil
	}

	return true, nil
}
