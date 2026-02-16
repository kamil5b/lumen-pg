package query_editor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (h *QueryEditorHandlerImplementation) HandleExecuteMultipleQueries(w http.ResponseWriter, r *http.Request) {
	// Get session from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Validate session
	session, err := h.authUC.ValidateSession(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get query
	query := r.FormValue("query")
	if strings.TrimSpace(query) == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='error'>Query cannot be empty</div>"))
		return
	}

	// Execute multiple queries
	results, err := h.queryUC.ExecuteMultipleQueries(r.Context(), session.Username, query)
	if err != nil {
		// Check for validation errors
		if validationErr, ok := err.(domain.ValidationError); ok {
			if validationErr.Field == "permission" {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("<div class='error'>" + validationErr.Message + "</div>"))
				return
			}
			if validationErr.Field == "query" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("<div class='error syntax error'>" + validationErr.Message + "</div>"))
				return
			}
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='error'>" + err.Error() + "</div>"))
		return
	}

	// Render results
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	var html strings.Builder
	html.WriteString("<div class='multiple-results'>")

	for i, result := range results {
		html.WriteString(fmt.Sprintf("<div class='result-set'><h3>Result %d</h3>", i+1))

		// Show row count
		if result.TotalCount > 0 {
			html.WriteString(fmt.Sprintf("<div class='info'>%d row(s) returned</div>", result.TotalCount))
		}

		// Render table if there are columns
		if len(result.Columns) > 0 {
			html.WriteString("<table class='query-results'><thead><tr>")
			for _, col := range result.Columns {
				html.WriteString("<th>" + col + "</th>")
			}
			html.WriteString("</tr></thead><tbody>")

			// Render rows
			for _, row := range result.Rows {
				html.WriteString("<tr>")
				for _, col := range result.Columns {
					value := row[col]
					var valueStr string
					if value == nil {
						valueStr = "NULL"
					} else {
						valueStr = fmt.Sprintf("%v", value)
					}
					html.WriteString("<td>" + valueStr + "</td>")
				}
				html.WriteString("</tr>")
			}

			html.WriteString("</tbody></table>")
		} else if result.RowCount > 0 {
			// DML query
			html.WriteString(fmt.Sprintf("<div class='success'>%d row(s) affected</div>", result.RowCount))
		} else {
			// DDL query or no results
			html.WriteString("<div class='success'>Query executed successfully</div>")
		}

		html.WriteString("</div>")
	}

	html.WriteString("</div>")
	w.Write([]byte(html.String()))
}
