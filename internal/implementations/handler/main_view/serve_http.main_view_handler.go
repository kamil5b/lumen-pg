package main_view

import "net/http"

func (h *MainViewHandlerImplementation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/main":
		h.HandleMainViewPage(w, r)
	case "/main/select-table":
		h.HandleTableSelect(w, r)
	case "/main/load-data":
		h.HandleLoadTableData(w, r)
	case "/main/filter":
		h.HandleFilterTable(w, r)
	case "/main/sort":
		h.HandleSortTable(w, r)
	case "/main/pagination/next":
		h.HandlePaginationNext(w, r)
	case "/main/pagination/previous":
		h.HandlePaginationPrevious(w, r)
	default:
		http.NotFound(w, r)
	}
}
