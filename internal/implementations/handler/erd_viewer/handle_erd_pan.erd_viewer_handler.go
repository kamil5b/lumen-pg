package erd_viewer

import (
	"net/http"
)

func (h *ERDViewerHandlerImplementation) HandleERDPan(w http.ResponseWriter, r *http.Request) {
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

	// Pan functionality is client-side (JavaScript), this endpoint acknowledges pan capability
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"pan_enabled": true, "draggable": true}`))
}
