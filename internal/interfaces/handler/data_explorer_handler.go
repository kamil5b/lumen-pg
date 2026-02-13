package handler

import "net/http"

// DataExplorerHandler handles data explorer sidebar HTTP requests
type DataExplorerHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	HandleLoadDataExplorer(w http.ResponseWriter, r *http.Request)
	HandleSelectDatabase(w http.ResponseWriter, r *http.Request)
	HandleSelectSchema(w http.ResponseWriter, r *http.Request)
	HandleSelectTable(w http.ResponseWriter, r *http.Request)
}
