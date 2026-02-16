package main_view

import (
	"net/http"
	"strconv"
)

func (h *MainViewHandlerImplementation) HandleSortTable(w http.ResponseWriter, r *http.Request) {
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
	column := r.FormValue("column")
	direction := r.FormValue("direction")

	if database == "" || schema == "" || table == "" || column == "" || direction == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Get pagination parameters
	offsetStr := r.FormValue("offset")
	offset := 0
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	limit := 50 // Default limit

	// Sort table data
	result, err := h.dataViewUC.SortTableData(r.Context(), session.Username, database, schema, table, column, direction, offset, limit)
	if err != nil {
		http.Error(w, "Error sorting table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render sorted table
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	html := `<table><thead><tr>`

	// Render column headers
	for _, col := range result.Columns {
		html += `<th>` + col + `</th>`
	}

	html += `</tr></thead><tbody>`

	// Render rows
	for _, row := range result.Rows {
		html += `<tr>`
		for _, col := range result.Columns {
			value := row[col]
			var valueStr string
			if value == nil {
				valueStr = "NULL"
			} else {
				valueStr = formatValue(value)
			}
			html += `<td>` + valueStr + `</td>`
		}
		html += `</tr>`
	}

	html += `</tbody></table>`

	w.Write([]byte(html))
}
