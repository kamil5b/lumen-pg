package handler

import "net/http"

// ERDViewerHandler handles ERD visualization HTTP requests
type ERDViewerHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	HandleERDViewerPage(w http.ResponseWriter, r *http.Request)
	HandleGenerateERD(w http.ResponseWriter, r *http.Request)
	HandleERDZoom(w http.ResponseWriter, r *http.Request)
	HandleERDPan(w http.ResponseWriter, r *http.Request)
	HandleTableClickInERD(w http.ResponseWriter, r *http.Request)
}
