package query_editor

import (
	"net/http"
)

func (h *QueryEditorHandlerImplementation) HandleQueryEditorPage(w http.ResponseWriter, r *http.Request) {
	// Get session from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Validate session
	_, err = h.authUC.ValidateSession(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Return query editor page HTML
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<title>Query Editor</title>
	<style>
		.query-editor { width: 100%; height: 300px; }
		.results-panel { margin-top: 20px; }
		.syntax-highlight.sql { font-family: monospace; }
	</style>
</head>
<body>
	<div class="query-editor-container">
		<h1>SQL Query Editor</h1>
		<form method="POST" action="/api/query/execute">
			<textarea name="query" class="query-editor syntax-highlight sql" placeholder="Enter your SQL query here..."></textarea>
			<button type="submit">Execute</button>
		</form>
		<div class="results-panel" id="results">
			<!-- Query results will be displayed here -->
		</div>
	</div>
</body>
</html>`))
}
