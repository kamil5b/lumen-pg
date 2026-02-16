package transaction

import (
	"net/http"
)

func (h *TransactionHandlerImplementation) HandleCommitTransaction(w http.ResponseWriter, r *http.Request) {
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

	// Check if transaction is expired
	isExpired, err := h.transactionUC.IsTransactionExpired(r.Context(), session.Username)
	if err != nil {
		http.Error(w, "Error checking transaction expiration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if isExpired {
		w.WriteHeader(http.StatusRequestTimeout)
		w.Write([]byte("<div class='error'>Transaction expired</div>"))
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

	// Commit transaction
	err = h.transactionUC.CommitTransaction(r.Context(), session.Username)
	if err != nil {
		http.Error(w, "Error committing transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<div class='success'>Transaction committed successfully</div>"))
}
