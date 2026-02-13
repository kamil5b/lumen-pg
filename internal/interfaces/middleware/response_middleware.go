package middleware

import "net/http"

// ResponseMiddleware transforms responses
type ResponseMiddleware interface {
	// ContentNegotiation handles content type negotiation
	ContentNegotiation(next http.Handler) http.Handler

	// CompressResponse compresses response body
	CompressResponse(next http.Handler) http.Handler

	// SetDefaultHeaders sets default response headers
	SetDefaultHeaders(next http.Handler) http.Handler
}
