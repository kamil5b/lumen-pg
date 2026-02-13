package middleware

import "net/http"

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
