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

	// Get active transaction
	txnState, err := h.transactionUC.GetActiveTransaction(r.Context(), session.Username)
	if err != nil {
		// Check if there are edits/deletes/inserts
		edits, _ := h.transactionUC.GetTransactionEdits(r.Context(), session.Username)
		deletes, _ := h.transactionUC.GetTransactionDeletes(r.Context(), session.Username)
		inserts, _ := h.transactionUC.GetTransactionInserts(r.Context(), session.Username)

		if len(edits) == 0 && len(deletes) == 0 && len(inserts) == 0 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("<div class='no-transaction'>No active transaction</div>"))
			return
		}

		// Return buffer information
		response := map[string]interface{}{
			"active":       false,
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

	// Get remaining time
	remainingTime, err := h.transactionUC.GetTransactionRemainingTime(r.Context(), session.Username)
	if err != nil {
		remainingTime = 0
	}

	// Get transaction edits
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

	// Build response
	response := map[string]interface{}{
		"active":         true,
		"transaction_id": txnState.ID,
		"username":       txnState.Username,
		"remaining_time": remainingTime,
		"edits":          edits,
		"deletes":        deletes,
		"inserts":        inserts,
		"edit_count":     len(edits),
		"delete_count":   len(deletes),
		"insert_count":   len(inserts),
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
