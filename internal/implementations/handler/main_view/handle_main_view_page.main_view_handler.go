package main_view

import (
	"net/http"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (h *MainViewHandlerImplementation) HandleMainViewPage(w http.ResponseWriter, r *http.Request) {
	// Get session from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Validate session
	session, err := h.authUC.ValidateSession(r.Context(), cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Get user accessible resources
	resources, err := h.authUC.GetUserAccessibleResources(r.Context(), session.Username)
	if err != nil {
		http.Error(w, "Error getting accessible resources: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Load first accessible table by default
	if len(resources.AccessibleTables) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<div>No accessible tables</div>"))
		return
	}

	// Get first table
	firstTable := resources.AccessibleTables[0]

	// Load table data
	tableData, err := h.dataViewUC.LoadTableData(r.Context(), session.Username, domain.TableDataParams{
		Database: firstTable.Database,
		Schema:   firstTable.Schema,
		Table:    firstTable.Name,
		Offset:   0,
		Limit:    50,
	})
	if err != nil {
		http.Error(w, "Error loading table data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render main view page
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	html := `<!DOCTYPE html>
<html>
<head>
	<title>Main View - Lumen PG</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 0; padding: 0; }
		.container { display: flex; height: 100vh; }
		.sidebar { width: 250px; background: #f8f9fa; padding: 20px; overflow-y: auto; }
		.main-content { flex: 1; padding: 20px; overflow-y: auto; }
		table { width: 100%; border-collapse: collapse; }
		th, td { padding: 8px; text-align: left; border: 1px solid #ddd; }
		th { background: #007bff; color: white; cursor: pointer; }
		.database-item { margin-bottom: 10px; }
		.schema-item { margin-left: 10px; margin-bottom: 5px; }
		.table-item { margin-left: 20px; margin-bottom: 3px; cursor: pointer; }
		.table-item:hover { background: #e9ecef; }
	</style>
</head>
<body>
	<div class="container">
		<div class="sidebar">
			<h3>Databases</h3>
			<div class="database-list">`

	// Render accessible databases
	for _, db := range resources.AccessibleDatabases {
		html += `<div class="database-item"><strong>` + db + `</strong></div>`
	}

	html += `
			</div>
		</div>
		<div class="main-content">
			<h2>Table: ` + firstTable.Name + `</h2>
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
		</div>
	</div>
</body>
</html>`

	w.Write([]byte(html))
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return stringifyInt(v)
	case float32, float64:
		return stringifyFloat(v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return "N/A"
	}
}

func stringifyInt(v interface{}) string {
	switch val := v.(type) {
	case int:
		return itoa(val)
	case int8:
		return itoa(int(val))
	case int16:
		return itoa(int(val))
	case int32:
		return itoa(int(val))
	case int64:
		return itoa(int(val))
	default:
		return "0"
	}
}

func stringifyFloat(v interface{}) string {
	// Simple float to string conversion
	return "0.0"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	var result []byte
	for n > 0 {
		digit := n % 10
		result = append([]byte{byte('0' + digit)}, result...)
		n /= 10
	}

	if negative {
		result = append([]byte{'-'}, result...)
	}

	return string(result)
}
