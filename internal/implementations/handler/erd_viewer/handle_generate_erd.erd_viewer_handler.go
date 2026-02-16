package erd_viewer

import (
	"encoding/json"
	"net/http"
)

func (h *ERDViewerHandlerImplementation) HandleGenerateERD(w http.ResponseWriter, r *http.Request) {
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

	// Get query parameters
	database := r.URL.Query().Get("database")
	schema := r.URL.Query().Get("schema")

	if database == "" || schema == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Generate ERD
	erdData, err := h.erdUC.GenerateERD(r.Context(), session.Username, database, schema)
	if err != nil {
		http.Error(w, "Error generating ERD: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(erdData)
}
