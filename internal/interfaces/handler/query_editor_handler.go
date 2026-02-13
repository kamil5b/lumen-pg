package handler

import "net/http"

// QueryEditorHandler handles query execution HTTP requests
type QueryEditorHandler interface {
	HTTPHandler
	HandleQueryEditorPage(w http.ResponseWriter, r *http.Request)
	HandleExecuteQuery(w http.ResponseWriter, r *http.Request)
	HandleExecuteMultipleQueries(w http.ResponseWriter, r *http.Request)
}
