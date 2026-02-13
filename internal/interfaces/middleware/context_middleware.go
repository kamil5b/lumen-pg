package middleware

import "net/http"

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
