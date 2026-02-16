package transaction

import (
	"net/http"
	"strconv"
)

func (h *TransactionHandlerImplementation) HandleDeleteRow(w http.ResponseWriter, r *http.Request) {
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
	rowIndexStr := r.FormValue("row_index")

	if database == "" || schema == "" || table == "" || rowIndexStr == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Convert row_index to int
	rowIndex, err := strconv.Atoi(rowIndexStr)
	if err != nil {
		http.Error(w, "Invalid row_index: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Delete row
	err = h.transactionUC.DeleteRow(r.Context(), session.Username, database, schema, table, rowIndex)
	if err != nil {
		http.Error(w, "Error deleting row: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<div class='success deleted'>Row marked for deletion</div>"))
}
