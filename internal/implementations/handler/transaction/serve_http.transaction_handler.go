package transaction

import "net/http"

func (h *TransactionHandlerImplementation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/transaction/start":
		h.HandleStartTransaction(w, r)
	case "/transaction/edit-cell":
		h.HandleEditCell(w, r)
	case "/transaction/delete-row":
		h.HandleDeleteRow(w, r)
	case "/transaction/insert-row":
		h.HandleInsertRow(w, r)
	case "/transaction/commit":
		h.HandleCommitTransaction(w, r)
	case "/transaction/rollback":
		h.HandleRollbackTransaction(w, r)
	case "/transaction/status":
		h.HandleGetTransactionStatus(w, r)
	default:
		http.NotFound(w, r)
	}
}
