package middleware

import "net/http"

// ValidationMiddleware validates request input
type ValidationMiddleware interface {
	// ValidateQueryParams validates query parameters
	ValidateQueryParams(next http.Handler) http.Handler

	// ValidateRequestBody validates request body
	ValidateRequestBody(next http.Handler) http.Handler

	// ValidateWhereClause validates SQL WHERE clause for injection
	ValidateWhereClause(next http.Handler) http.Handler

	// ValidateSQLQuery validates SQL query for injection
	ValidateSQLQuery(next http.Handler) http.Handler
}
