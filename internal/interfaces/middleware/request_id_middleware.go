package middleware

import "net/http"

// RequestIDMiddleware generates and injects request IDs
type RequestIDMiddleware interface {
	// InjectRequestID generates and injects a request ID
	InjectRequestID(next http.Handler) http.Handler
}
