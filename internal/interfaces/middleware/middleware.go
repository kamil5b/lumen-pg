package middleware

import "net/http"

// Middleware represents an HTTP middleware function
type Middleware func(http.Handler) http.Handler

// AuthenticationMiddleware validates user sessions before allowing access
type AuthenticationMiddleware interface {
	// Authenticate checks if the request has a valid session and injects user context
	Authenticate(next http.Handler) http.Handler

	// RequireAuth is a middleware that requires authentication
	RequireAuth(next http.Handler) http.Handler

	// OptionalAuth is a middleware that optionally injects user context if authenticated
	OptionalAuth(next http.Handler) http.Handler
}

// AuthorizationMiddleware validates user permissions for resources
type AuthorizationMiddleware interface {
	// RequireTableAccess checks if the user has access to a table
	RequireTableAccess(next http.Handler) http.Handler

	// RequireSelectPermission checks if the user can SELECT from a table
	RequireSelectPermission(next http.Handler) http.Handler

	// RequireInsertPermission checks if the user can INSERT into a table
	RequireInsertPermission(next http.Handler) http.Handler

	// RequireUpdatePermission checks if the user can UPDATE a table
	RequireUpdatePermission(next http.Handler) http.Handler

	// RequireDeletePermission checks if the user can DELETE from a table
	RequireDeletePermission(next http.Handler) http.Handler

	// RequireDatabaseAccess checks if the user can access a database
	RequireDatabaseAccess(next http.Handler) http.Handler
}

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

// SecurityMiddleware handles security concerns
type SecurityMiddleware interface {
	// SetSecurityHeaders sets appropriate security headers
	SetSecurityHeaders(next http.Handler) http.Handler

	// PreventCSRF prevents cross-site request forgery
	PreventCSRF(next http.Handler) http.Handler

	// RateLimiter limits request rate per user/IP
	RateLimiter(next http.Handler) http.Handler

	// ValidateCookieIntegrity validates that cookies haven't been tampered with
	ValidateCookieIntegrity(next http.Handler) http.Handler

	// EnforceSameSiteCookie enforces SameSite cookie attribute
	EnforceSameSiteCookie(next http.Handler) http.Handler

	// RequireHTTPSForCookies enforces HTTPS for secure cookies
	RequireHTTPSForCookies(next http.Handler) http.Handler
}

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

// ContextMiddleware manages request context
type ContextMiddleware interface {
	// InjectUser injects user information into the request context
	InjectUser(next http.Handler) http.Handler

	// InjectSession injects session information into the request context
	InjectSession(next http.Handler) http.Handler

	// InjectTransaction injects active transaction information into the request context
	InjectTransaction(next http.Handler) http.Handler

	// InjectUserPermissions injects user permissions into the request context
	InjectUserPermissions(next http.Handler) http.Handler

	// InjectMetadata injects cached metadata into the request context
	InjectMetadata(next http.Handler) http.Handler
}

// ErrorHandlingMiddleware handles errors and panics
type ErrorHandlingMiddleware interface {
	// HandleErrors catches and processes errors
	HandleErrors(next http.Handler) http.Handler

	// RecoverFromPanic recovers from panics
	RecoverFromPanic(next http.Handler) http.Handler

	// ValidateHTTPMethod validates HTTP method
	ValidateHTTPMethod(allowedMethods ...string) func(http.Handler) http.Handler
}

// ResponseMiddleware transforms responses
type ResponseMiddleware interface {
	// ContentNegotiation handles content type negotiation
	ContentNegotiation(next http.Handler) http.Handler

	// CompressResponse compresses response body
	CompressResponse(next http.Handler) http.Handler

	// SetDefaultHeaders sets default response headers
	SetDefaultHeaders(next http.Handler) http.Handler
}

// CORSMiddleware handles cross-origin requests
type CORSMiddleware interface {
	// HandleCORS handles cross-origin resource sharing
	HandleCORS(next http.Handler) http.Handler
}

// RequestIDMiddleware generates and injects request IDs
type RequestIDMiddleware interface {
	// InjectRequestID generates and injects a request ID
	InjectRequestID(next http.Handler) http.Handler
}
