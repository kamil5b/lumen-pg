package middleware

import "net/http"

// LoggingMiddleware handles request/response logging
type LoggingMiddleware interface {
	// LogRequest logs incoming requests
	LogRequest(next http.Handler) http.Handler

	// LogQueryExecution logs SQL query executions
	LogQueryExecution(next http.Handler) http.Handler

	// LogSecurityEvents logs security-related events
	LogSecurityEvents(next http.Handler) http.Handler

	// LogTransactionEvents logs transaction-related events
	LogTransactionEvents(next http.Handler) http.Handler
}
