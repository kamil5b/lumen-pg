package erd_viewer

import "net/http"

func (h *ERDViewerHandlerImplementation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/erd":
		h.HandleERDViewerPage(w, r)
	case "/erd/generate":
		h.HandleGenerateERD(w, r)
	case "/erd/zoom":
		h.HandleERDZoom(w, r)
	case "/erd/pan":
		h.HandleERDPan(w, r)
	case "/erd/table":
		h.HandleTableClickInERD(w, r)
	default:
		http.NotFound(w, r)
	}
}
