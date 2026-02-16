package transaction

import (
	"encoding/json"
	"net/http"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (h *TransactionHandlerImplementation) HandleGetTransactionStatus(w http.ResponseWriter, r *http.Request) {
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

	// Get transaction edits first
	edits, err := h.transactionUC.GetTransactionEdits(r.Context(), session.Username)
	if err != nil {
		edits = make(map[int]domain.RowEdit)
	}

	// Get transaction deletes
	deletes, err := h.transactionUC.GetTransactionDeletes(r.Context(), session.Username)
	if err != nil {
		deletes = []int{}
	}

	// Get transaction inserts
	inserts, err := h.transactionUC.GetTransactionInserts(r.Context(), session.Username)
	if err != nil {
		inserts = []domain.RowInsert{}
	}

	// Check if there are any edits/deletes/inserts (buffer data)
	if len(edits) > 0 || len(deletes) > 0 || len(inserts) > 0 {
		// Return buffer information (may or may not have active transaction)
		response := map[string]interface{}{
			"edits":        edits,
			"deletes":      deletes,
			"inserts":      inserts,
			"edit_count":   len(edits),
			"delete_count": len(deletes),
			"insert_count": len(inserts),
		}

		jsonData, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
		return
	}

	// No buffer data - check for active transaction
	txnState, err := h.transactionUC.GetActiveTransaction(r.Context(), session.Username)
	if err != nil {
		// No active transaction and no buffer data
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<div class='no-transaction'>No active transaction</div>"))
		return
	}

	// Get remaining time
	remainingTime, err := h.transactionUC.GetTransactionRemainingTime(r.Context(), session.Username)
	if err != nil {
		remainingTime = 0
	}

	// Build response for active transaction
	response := map[string]interface{}{
		"active":         true,
		"transaction_id": txnState.ID,
		"username":       txnState.Username,
		"remaining_time": remainingTime,
		"edits":          edits,
		"deletes":        deletes,
		"inserts":        inserts,
		"edit_count":     0,
		"delete_count":   0,
		"insert_count":   0,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshaling response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
