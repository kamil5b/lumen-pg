package handler

import "net/http"

// HTTPHandler defines common HTTP handler interface
type HTTPHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// LoginHandler handles authentication-related HTTP requests
type LoginHandler interface {
	HTTPHandler
	HandleLoginPage(w http.ResponseWriter, r *http.Request)
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleLogout(w http.ResponseWriter, r *http.Request)
}

// MainViewHandler handles main data view HTTP requests
type MainViewHandler interface {
	HTTPHandler
	HandleMainViewPage(w http.ResponseWriter, r *http.Request)
	HandleTableSelect(w http.ResponseWriter, r *http.Request)
	HandleLoadTableData(w http.ResponseWriter, r *http.Request)
	HandleFilterTable(w http.ResponseWriter, r *http.Request)
	HandleSortTable(w http.ResponseWriter, r *http.Request)
	HandlePaginationNext(w http.ResponseWriter, r *http.Request)
	HandlePaginationPrevious(w http.ResponseWriter, r *http.Request)
}

// QueryEditorHandler handles query execution HTTP requests
type QueryEditorHandler interface {
	HTTPHandler
	HandleQueryEditorPage(w http.ResponseWriter, r *http.Request)
	HandleExecuteQuery(w http.ResponseWriter, r *http.Request)
	HandleExecuteMultipleQueries(w http.ResponseWriter, r *http.Request)
}

// ERDViewerHandler handles ERD visualization HTTP requests
type ERDViewerHandler interface {
	HTTPHandler
	HandleERDViewerPage(w http.ResponseWriter, r *http.Request)
	HandleGenerateERD(w http.ResponseWriter, r *http.Request)
	HandleERDZoom(w http.ResponseWriter, r *http.Request)
	HandleERDPan(w http.ResponseWriter, r *http.Request)
	HandleTableClickInERD(w http.ResponseWriter, r *http.Request)
}

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

// DataExplorerHandler handles data explorer sidebar HTTP requests
type DataExplorerHandler interface {
	HTTPHandler
	HandleLoadDataExplorer(w http.ResponseWriter, r *http.Request)
	HandleSelectDatabase(w http.ResponseWriter, r *http.Request)
	HandleSelectSchema(w http.ResponseWriter, r *http.Request)
	HandleSelectTable(w http.ResponseWriter, r *http.Request)
}

// NavigationHandler handles navigation-related HTTP requests
type NavigationHandler interface {
	HTTPHandler
	HandleNavigateToParentRow(w http.ResponseWriter, r *http.Request)
	HandleNavigateToChildRows(w http.ResponseWriter, r *http.Request)
	HandleGetChildTableReferences(w http.ResponseWriter, r *http.Request)
}

// MetadataHandler handles metadata-related HTTP requests
type MetadataHandler interface {
	HTTPHandler
	HandleRefreshMetadata(w http.ResponseWriter, r *http.Request)
	HandleGetDatabaseMetadata(w http.ResponseWriter, r *http.Request)
	HandleGetSchemaMetadata(w http.ResponseWriter, r *http.Request)
	HandleGetTableMetadata(w http.ResponseWriter, r *http.Request)
}

// HealthHandler handles health check HTTP requests
type HealthHandler interface {
	HTTPHandler
	HandleHealthCheck(w http.ResponseWriter, r *http.Request)
	HandleIsInitialized(w http.ResponseWriter, r *http.Request)
}
