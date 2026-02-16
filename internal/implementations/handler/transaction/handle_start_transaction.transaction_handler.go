package transaction

import (
	"net/http"
)

func (h *TransactionHandlerImplementation) HandleStartTransaction(w http.ResponseWriter, r *http.Request) {
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

	// Check update permission
	hasPermission, err := h.rbacUC.CheckUpdatePermission(r.Context(), session.Username, database, schema, table)
	if err != nil {
		http.Error(w, "Error checking permissions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !hasPermission {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("<div class='error'>permission denied</div>"))
		return
	}

	// Check if transaction already active
	hasActive, err := h.transactionUC.CheckActiveTransaction(r.Context(), session.Username)
	if err != nil {
		http.Error(w, "Error checking active transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if hasActive {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("<div class='error'>Transaction already active</div>"))
		return
	}

	// Start transaction
	txnState, err := h.transactionUC.StartTransaction(r.Context(), session.Username, database, schema, table)
	if err != nil {
		http.Error(w, "Error starting transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response with transaction indicator
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<div class='transaction-active'>
		<span class='indicator'>Transaction Active</span>
		<span class='timer' data-seconds='60'>60s</span>
		<span class='txn-id' style='display:none'>` + txnState.ID + `</span>
	</div>`))
}
