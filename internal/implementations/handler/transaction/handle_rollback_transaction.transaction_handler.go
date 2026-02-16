package transaction

import (
	"net/http"
)

func (h *TransactionHandlerImplementation) HandleRollbackTransaction(w http.ResponseWriter, r *http.Request) {
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

	// Rollback transaction
	err = h.transactionUC.RollbackTransaction(r.Context(), session.Username)
	if err != nil {
		http.Error(w, "Error rolling back transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<div class='success'>Transaction rolled back successfully</div>"))
}
