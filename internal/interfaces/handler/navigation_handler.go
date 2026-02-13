package handler

import "net/http"

// NavigationHandler handles navigation-related HTTP requests
type NavigationHandler interface {
	HTTPHandler
	HandleNavigateToParentRow(w http.ResponseWriter, r *http.Request)
	HandleNavigateToChildRows(w http.ResponseWriter, r *http.Request)
	HandleGetChildTableReferences(w http.ResponseWriter, r *http.Request)
}
