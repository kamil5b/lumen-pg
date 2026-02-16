package main_view

import (
	"fmt"
	"net/http"
	"strconv"
)

func (h *MainViewHandlerImplementation) HandleFilterTable(w http.ResponseWriter, r *http.Request) {
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

	// Get parameters
	database := r.FormValue("database")
	schema := r.FormValue("schema")
	table := r.FormValue("table")
	whereClause := r.FormValue("where")

	if database == "" || schema == "" || table == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Validate WHERE clause
	isValid, err := h.dataViewUC.ValidateWhereClause(r.Context(), whereClause)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='error'>" + err.Error() + "</div>"))
		return
	}
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='error'>Invalid WHERE clause</div>"))
		return
	}

	// Get pagination parameters
	offsetStr := r.FormValue("offset")
	offset := 0
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	limit := 50 // Default limit

	// Filter table data
	result, err := h.dataViewUC.FilterTableData(r.Context(), session.Username, database, schema, table, whereClause, offset, limit)
	if err != nil {
		http.Error(w, "Error filtering table data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render filtered results
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	html := `<table class="filtered-results">
		<thead>
			<tr>`

	for _, col := range result.Columns {
		html += `<th>` + col + `</th>`
	}

	html += `</tr>
		</thead>
		<tbody>`

	for _, row := range result.Rows {
		html += `<tr>`
		for _, col := range result.Columns {
			value := row[col]
			var valueStr string
			if value == nil {
				valueStr = "NULL"
			} else {
				valueStr = fmt.Sprintf("%v", value)
			}
			html += `<td>` + valueStr + `</td>`
		}
		html += `</tr>`
	}

	html += `</tbody>
	</table>`

	w.Write([]byte(html))
}
