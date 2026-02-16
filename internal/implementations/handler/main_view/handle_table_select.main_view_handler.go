package main_view

import (
	"net/http"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (h *MainViewHandlerImplementation) HandleTableSelect(w http.ResponseWriter, r *http.Request) {
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

	if database == "" || schema == "" || table == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Check table access
	hasAccess, err := h.rbacUC.CheckTableAccess(r.Context(), session.Username, database, schema, table)
	if err != nil {
		http.Error(w, "Error checking table access: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !hasAccess {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("<div class='error'>Access denied to table</div>"))
		return
	}

	// Load table data
	tableData, err := h.dataViewUC.LoadTableData(r.Context(), session.Username, domain.TableDataParams{
		Database: database,
		Schema:   schema,
		Table:    table,
		Offset:   0,
		Limit:    50,
	})
	if err != nil {
		http.Error(w, "Error loading table data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render table data as HTML
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	html := `<div class="table-data">
		<h2>Table: ` + table + `</h2>
		<table>
			<thead>
				<tr>`

	// Render column headers
	for _, col := range tableData.Columns {
		html += `<th>` + col + `</th>`
	}

	html += `
				</tr>
			</thead>
			<tbody>`

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

	html += `
			</tbody>
		</table>
	</div>`

	w.Write([]byte(html))
}
