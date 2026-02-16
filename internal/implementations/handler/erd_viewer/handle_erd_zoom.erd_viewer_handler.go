package erd_viewer

import (
	"net/http"
)

func (h *ERDViewerHandlerImplementation) HandleERDZoom(w http.ResponseWriter, r *http.Request) {
	// Get session from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Validate session
	_, err = h.authUC.ValidateSession(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get zoom parameter
	zoomLevel := r.URL.Query().Get("level")
	if zoomLevel == "" {
		http.Error(w, "Missing zoom level parameter", http.StatusBadRequest)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<div class='zoom-success' data-level='` + zoomLevel + `'>Zoom applied</div>`))
}
