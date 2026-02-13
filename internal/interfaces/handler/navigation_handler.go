package handler

import "net/http"

// NavigationHandler handles navigation-related HTTP requests
type NavigationHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	HandleNavigateToParentRow(w http.ResponseWriter, r *http.Request)
	HandleNavigateToChildRows(w http.ResponseWriter, r *http.Request)
	HandleGetChildTableReferences(w http.ResponseWriter, r *http.Request)
}
