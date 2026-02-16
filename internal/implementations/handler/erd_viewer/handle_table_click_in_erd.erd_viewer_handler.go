package erd_viewer

import (
	"net/http"
)

func (h *ERDViewerHandlerImplementation) HandleTableClickInERD(w http.ResponseWriter, r *http.Request) {
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

	if database == "" || schema == "" || table == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Get table metadata
	tableData, err := h.erdUC.GetTableBoxData(r.Context(), session.Username, database, schema, table)
	if err != nil {
		http.Error(w, "Error getting table metadata: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render table details panel
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	html := `<div class="table-details-panel">
		<h3>Table: ` + table + `</h3>
		<table class="table-metadata">
			<thead>
				<tr>
					<th>Column</th>
					<th>Type</th>
					<th>Nullable</th>
					<th>Key</th>
				</tr>
			</thead>
			<tbody>`

	for _, col := range tableData.Columns {
		nullable := "NO"
		if col.IsNullable {
			nullable = "YES"
		}

		keyType := ""
		if col.IsPrimary {
			keyType = "PK"
		}

		// Check if foreign key
		for _, fk := range tableData.ForeignKeys {
			if fk.ColumnName == col.Name {
				if keyType != "" {
					keyType += ", FK"
				} else {
					keyType = "FK"
				}
				break
			}
		}

		html += `<tr>
			<td>` + col.Name + `</td>
			<td>` + col.DataType + `</td>
			<td>` + nullable + `</td>
			<td>` + keyType + `</td>
		</tr>`
	}

	html += `</tbody>
		</table>`

	// Show foreign keys if any
	if len(tableData.ForeignKeys) > 0 {
		html += `<h4>Foreign Keys</h4><ul>`
		for _, fk := range tableData.ForeignKeys {
			html += `<li>` + fk.ColumnName + ` → ` + fk.ReferencedTable + `.` + fk.ReferencedColumn + `</li>`
		}
		html += `</ul>`
	}

	html += `</div>`

	w.Write([]byte(html))
}
