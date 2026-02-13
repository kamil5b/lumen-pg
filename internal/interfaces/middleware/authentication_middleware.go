package middleware

import "net/http"

// AuthenticationMiddleware validates user sessions before allowing access
type AuthenticationMiddleware interface {
	// Authenticate checks if the request has a valid session and injects user context
	Authenticate(next http.Handler) http.Handler

	// RequireAuth is a middleware that requires authentication
	RequireAuth(next http.Handler) http.Handler

	// OptionalAuth is a middleware that optionally injects user context if authenticated
	OptionalAuth(next http.Handler) http.Handler
}
