package transaction

import (
	"net/http"
)

func (h *TransactionHandlerImplementation) HandleInsertRow(w http.ResponseWriter, r *http.Request) {
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

	// Check if transaction is active
	hasActive, err := h.transactionUC.CheckActiveTransaction(r.Context(), session.Username)
	if err != nil {
		http.Error(w, "Error checking active transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !hasActive {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='error'>No active transaction</div>"))
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get parameters
	database := r.URL.Query().Get("database")
	schema := r.URL.Query().Get("schema")
	table := r.URL.Query().Get("table")

	if database == "" || schema == "" || table == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Build row data from form values
	rowData := make(map[string]interface{})
	for key, values := range r.Form {
		if len(values) > 0 {
			rowData[key] = values[0]
		}
	}

	// Insert row
	err = h.transactionUC.InsertRow(r.Context(), session.Username, database, schema, table, rowData)
	if err != nil {
		http.Error(w, "Error inserting row: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<div class='success inserted'>New row added to buffer</div>"))
}
