package middleware

import "net/http"

// Middleware represents an HTTP middleware function
type Middleware func(http.Handler) http.Handler
