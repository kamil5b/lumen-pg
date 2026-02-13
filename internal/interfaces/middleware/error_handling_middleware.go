package middleware

import "net/http"

// ErrorHandlingMiddleware handles errors and panics
type ErrorHandlingMiddleware interface {
	// HandleErrors catches and processes errors
	HandleErrors(next http.Handler) http.Handler

	// RecoverFromPanic recovers from panics
	RecoverFromPanic(next http.Handler) http.Handler

	// ValidateHTTPMethod validates HTTP method
	ValidateHTTPMethod(allowedMethods ...string) func(http.Handler) http.Handler
}
