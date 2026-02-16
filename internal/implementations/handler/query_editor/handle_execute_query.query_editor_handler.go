package query_editor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (h *QueryEditorHandlerImplementation) HandleExecuteQuery(w http.ResponseWriter, r *http.Request) {
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

	// Get pagination parameters
	offsetStr := r.FormValue("offset")
	offset := 0
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	limit := 50 // Default limit

	// Check for pagination request
	if offsetStr != "" {
		// Execute with pagination
		result, err := h.queryUC.ExecuteQueryWithPagination(r.Context(), session.Username, domain.QueryParams{
			Query:  query,
			Offset: offset,
			Limit:  limit,
		})
		if err != nil {
			// Check for permission error
			if validationErr, ok := err.(domain.ValidationError); ok {
				if validationErr.Field == "permission" {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte("<div class='error'>" + validationErr.Message + "</div>"))
					return
				}
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("<div class='error'>" + err.Error() + "</div>"))
			return
		}

		// Render result with pagination info
		h.renderQueryResult(w, result)
		return
	}

	// Execute query
	result, err := h.queryUC.ExecuteQuery(r.Context(), session.Username, query, offset, limit)
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

	// Render result
	h.renderQueryResult(w, result)
}

func (h *QueryEditorHandlerImplementation) renderQueryResult(w http.ResponseWriter, result *domain.QueryResult) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	var html strings.Builder

	// Show pagination info if applicable
	if result.TotalCount > 1000 {
		html.WriteString(fmt.Sprintf("<div class='pagination-info'>Data size: %d rows (only first 1000 are accessible)</div>", result.TotalCount))
	} else if result.TotalCount > 0 {
		html.WriteString(fmt.Sprintf("<div class='pagination-info'>Data size: %d rows</div>", result.TotalCount))
	}

	// Check for hard limit
	if result.RowCount == 0 && result.TotalCount > 1000 {
		html.WriteString("<div class='warning'>hard limit of 1000 rows reached</div>")
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
					// Convert value to string
					switch v := value.(type) {
					case string:
						valueStr = v
					case int, int8, int16, int32, int64:
						valueStr = fmt.Sprintf("%d", v)
					case float32, float64:
						valueStr = fmt.Sprintf("%f", v)
					case bool:
						valueStr = fmt.Sprintf("%t", v)
					default:
						// Try JSON encoding for complex types
						jsonBytes, err := json.Marshal(v)
						if err == nil {
							valueStr = string(jsonBytes)
						} else {
							valueStr = fmt.Sprintf("%v", v)
						}
					}
				}
				html.WriteString("<td>" + valueStr + "</td>")
			}
			html.WriteString("</tr>")
		}

		html.WriteString("</tbody></table>")

		// Add pagination controls if needed
		if result.TotalCount > result.RowCount {
			html.WriteString("<div class='pagination'>")
			html.WriteString(fmt.Sprintf("<span>Showing %d of %d rows</span>", result.RowCount, result.TotalCount))
			html.WriteString("</div>")
		}
	} else if result.RowCount > 0 {
		// DML query (INSERT, UPDATE, DELETE)
		html.WriteString(fmt.Sprintf("<div class='success'>Query executed successfully. %d row(s) affected.</div>", result.RowCount))
	} else {
		// DDL query or no results
		html.WriteString("<div class='success'>Query executed successfully.</div>")
	}

	w.Write([]byte(html.String()))
}
