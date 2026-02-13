package handler

import "net/http"

// HTTPHandler defines common HTTP handler interface
type HTTPHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// HealthHandler handles health check HTTP requests
type HealthHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	HandleHealthCheck(w http.ResponseWriter, r *http.Request)
	HandleIsInitialized(w http.ResponseWriter, r *http.Request)
}
