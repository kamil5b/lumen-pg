package middleware

import "net/http"

// CORSMiddleware handles cross-origin requests
type CORSMiddleware interface {
	// HandleCORS handles cross-origin resource sharing
	HandleCORS(next http.Handler) http.Handler
}
