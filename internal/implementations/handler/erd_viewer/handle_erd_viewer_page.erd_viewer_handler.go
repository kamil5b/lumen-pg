package erd_viewer

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (h *ERDViewerHandlerImplementation) HandleERDViewerPage(w http.ResponseWriter, r *http.Request) {
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

	if database == "" {
		http.Error(w, "Missing database parameter", http.StatusBadRequest)
		return
	}

	// Get available schemas
	schemas, err := h.erdUC.GetAvailableSchemas(r.Context(), session.Username, database)
	if err != nil {
		http.Error(w, "Error getting schemas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Default to first schema if not specified
	if schema == "" && len(schemas) > 0 {
		schema = schemas[0]
	}

	// Check if schema is empty
	isEmpty, err := h.erdUC.IsSchemaEmpty(r.Context(), session.Username, database, schema)
	if err == nil && isEmpty {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		h.renderEmptySchemaPage(w, database, schema, schemas)
		return
	}

	// Generate ERD data
	erdDataRaw, err := h.erdUC.GenerateERD(r.Context(), session.Username, database, schema)
	if err != nil {
		http.Error(w, "Error generating ERD: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Type assert to domain.ERDData
	erdData, ok := erdDataRaw.(*domain.ERDData)
	if !ok {
		http.Error(w, "Invalid ERD data type", http.StatusInternalServerError)
		return
	}

	// Render ERD page
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	h.renderERDPage(w, database, schema, schemas, erdData)
}

func (h *ERDViewerHandlerImplementation) renderERDPage(w http.ResponseWriter, database, schema string, schemas []string, erdData *domain.ERDData) {
	var html strings.Builder

	html.WriteString(`<!DOCTYPE html>
<html>
<head>
	<title>ERD Viewer - ` + schema + `</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
		.header { margin-bottom: 20px; }
		.schema-selector { margin-bottom: 20px; }
		.erd-container { position: relative; width: 100%; height: 600px; border: 1px solid #ccc; overflow: hidden; }
		.erd-canvas { width: 100%; height: 100%; position: relative; cursor: move; }
		.erd-diagram { position: relative; }
		.erd-table { position: absolute; border: 1px solid #333; background: white; padding: 10px; min-width: 150px; }
		.table-name { font-weight: bold; margin-bottom: 5px; background: #f0f0f0; padding: 5px; }
		.table-column { padding: 3px; font-size: 12px; }
		.controls { margin-bottom: 10px; }
		.zoom-controls button { margin-right: 5px; }
		.draggable { cursor: move; }
		.pan { cursor: grab; }
	</style>
</head>
<body>
	<div class="header">
		<h1>ERD Viewer</h1>
		<div class="schema-selector">
			<label>Database: ` + database + `</label>
			<label>Schema:</label>
			<select id="schema-select" onchange="location.href='/erd?database=` + database + `&schema='+this.value">`)

	for _, s := range schemas {
		selected := ""
		if s == schema {
			selected = " selected"
		}
		html.WriteString(fmt.Sprintf(`<option value="%s"%s>%s</option>`, s, selected, s))
	}

	html.WriteString(`
			</select>
		</div>
	</div>
	<div class="controls">
		<div class="zoom-controls">
			<button id="zoom-in" class="zoom-in">Zoom In</button>
			<button id="zoom-out" class="zoom-out">Zoom Out</button>
			<button id="zoom-reset" class="zoom-reset">Reset</button>
		</div>
	</div>
	<div class="erd-container">
		<div class="erd-canvas draggable pan" id="erd-canvas">
			<div class="erd-diagram" id="erd-diagram">`)

	// Render ERD tables
	if erdData != nil {
		for i, table := range erdData.Tables {
			html.WriteString(fmt.Sprintf(`
				<div class="erd-table" style="left: %dpx; top: %dpx;">
					<div class="table-name">%s</div>`, 50+(i*200), 50+(i*100), table.Name))
			for _, col := range table.Columns {
				keyIndicator := ""
				if col.IsPrimary {
					keyIndicator = " (PK)"
				}
				html.WriteString(fmt.Sprintf(`
					<div class="table-column">%s: %s%s</div>`, col.Name, col.DataType, keyIndicator))
			}
			html.WriteString(`
				</div>`)
		}
	}

	html.WriteString(`
			</div>
		</div>
	</div>
	<div id="table-details" style="margin-top: 20px;">
		<!-- Table details will appear here when a table is clicked -->
	</div>
	<script>
		// Add interactivity for zoom and pan
		const canvas = document.getElementById('erd-canvas');
		let scale = 1;
		let isDragging = false;
		let startX, startY, scrollLeft, scrollTop;

		document.getElementById('zoom-in').addEventListener('click', () => {
			scale = Math.min(scale + 0.1, 2);
			document.getElementById('erd-diagram').style.transform = 'scale(' + scale + ')';
		});

		document.getElementById('zoom-out').addEventListener('click', () => {
			scale = Math.max(scale - 0.1, 0.5);
			document.getElementById('erd-diagram').style.transform = 'scale(' + scale + ')';
		});

		document.getElementById('zoom-reset').addEventListener('click', () => {
			scale = 1;
			document.getElementById('erd-diagram').style.transform = 'scale(1)';
		});

		// Pan functionality
		canvas.addEventListener('mousedown', (e) => {
			isDragging = true;
			startX = e.pageX - canvas.offsetLeft;
			startY = e.pageY - canvas.offsetTop;
			scrollLeft = canvas.scrollLeft;
			scrollTop = canvas.scrollTop;
		});

		canvas.addEventListener('mousemove', (e) => {
			if (!isDragging) return;
			e.preventDefault();
			const x = e.pageX - canvas.offsetLeft;
			const y = e.pageY - canvas.offsetTop;
			const walkX = (x - startX) * 1;
			const walkY = (y - startY) * 1;
			canvas.scrollLeft = scrollLeft - walkX;
			canvas.scrollTop = scrollTop - walkY;
		});

		canvas.addEventListener('mouseup', () => {
			isDragging = false;
		});

		canvas.addEventListener('mouseleave', () => {
			isDragging = false;
		});
	</script>
</body>
</html>`)

	w.Write([]byte(html.String()))
}

func (h *ERDViewerHandlerImplementation) renderEmptySchemaPage(w http.ResponseWriter, database, schema string, schemas []string) {
	var html strings.Builder

	html.WriteString(`<!DOCTYPE html>
<html>
<head>
	<title>ERD Viewer - ` + schema + `</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
		.empty-message { text-align: center; margin-top: 50px; color: #666; }
	</style>
</head>
<body>
	<h1>ERD Viewer</h1>
	<div class="schema-selector">
		<label>Database: ` + database + `</label>
		<label>Schema:</label>
		<select onchange="location.href='/erd?database=` + database + `&schema='+this.value">`)

	for _, s := range schemas {
		selected := ""
		if s == schema {
			selected = " selected"
		}
		html.WriteString(fmt.Sprintf(`<option value="%s"%s>%s</option>`, s, selected, s))
	}

	html.WriteString(`
		</select>
	</div>
	<div class="empty-message">
		<h2>No tables found in schema "` + schema + `"</h2>
	</div>
</body>
</html>`)

	w.Write([]byte(html.String()))
}
