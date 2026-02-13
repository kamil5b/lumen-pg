package handler

import "net/http"

// MetadataHandler handles metadata-related HTTP requests
type MetadataHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	HandleRefreshMetadata(w http.ResponseWriter, r *http.Request)
	HandleGetDatabaseMetadata(w http.ResponseWriter, r *http.Request)
	HandleGetSchemaMetadata(w http.ResponseWriter, r *http.Request)
	HandleGetTableMetadata(w http.ResponseWriter, r *http.Request)
}
