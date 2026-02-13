package handler

import "net/http"

// TransactionHandler handles transaction-related HTTP requests
type TransactionHandler interface {
	HTTPHandler
	HandleStartTransaction(w http.ResponseWriter, r *http.Request)
	HandleEditCell(w http.ResponseWriter, r *http.Request)
	HandleDeleteRow(w http.ResponseWriter, r *http.Request)
	HandleInsertRow(w http.ResponseWriter, r *http.Request)
	HandleCommitTransaction(w http.ResponseWriter, r *http.Request)
	HandleRollbackTransaction(w http.ResponseWriter, r *http.Request)
	HandleGetTransactionStatus(w http.ResponseWriter, r *http.Request)
}
