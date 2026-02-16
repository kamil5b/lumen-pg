package query_editor

import "net/http"

func (h *QueryEditorHandlerImplementation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/query-editor":
		h.HandleQueryEditorPage(w, r)
	case "/api/query/execute":
		h.HandleExecuteQuery(w, r)
	case "/api/query/execute-multiple":
		h.HandleExecuteMultipleQueries(w, r)
	default:
		http.NotFound(w, r)
	}
}
