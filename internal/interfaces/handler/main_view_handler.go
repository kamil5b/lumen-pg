package handler

import "net/http"

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
