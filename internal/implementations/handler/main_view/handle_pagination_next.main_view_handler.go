package main_view

import (
	"net/http"
	"strconv"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (h *MainViewHandlerImplementation) HandlePaginationNext(w http.ResponseWriter, r *http.Request) {
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

	// Get query parameters
	database := r.URL.Query().Get("database")
	schema := r.URL.Query().Get("schema")
	table := r.URL.Query().Get("table")
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	if database == "" || schema == "" || table == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Parse offset and limit
	offset := 0
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	limit := 50 // Default limit
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	// Move to next page
	offset += limit

	// Load table data
	tableData, err := h.dataViewUC.LoadTableData(r.Context(), session.Username, domain.TableDataParams{
		Database: database,
		Schema:   schema,
		Table:    table,
		Offset:   offset,
		Limit:    limit,
	})
	if err != nil {
		http.Error(w, "Error loading table data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render table data
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
