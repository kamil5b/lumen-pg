package main_view

import (
	"net/http"
	"strconv"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (h *MainViewHandlerImplementation) HandlePaginationPrevious(w http.ResponseWriter, r *http.Request) {
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
	offsetStr := r.FormValue("offset")

	if database == "" || schema == "" || table == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Parse offset
	offset := 0
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	// Calculate previous offset
	newOffset := offset - 50
	if newOffset < 0 {
		newOffset = 0
	}

	// Load table data with new offset
	tableData, err := h.dataViewUC.LoadTableData(r.Context(), session.Username, domain.TableDataParams{
		Database: database,
		Schema:   schema,
		Table:    table,
		Offset:   newOffset,
		Limit:    50,
	})
	if err != nil {
		http.Error(w, "Error loading table data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render table HTML
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	html := `<table><thead><tr>`

	// Render column headers
	for _, col := range tableData.Columns {
		html += `<th>` + col + `</th>`
	}

	html += `</tr></thead><tbody>`

	// Render rows
	for _, row := range tableData.Rows {
		html += `<tr>`
		for _, col := range tableData.Columns {
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
