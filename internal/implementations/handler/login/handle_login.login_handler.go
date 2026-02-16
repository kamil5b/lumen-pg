package login

import (
	"net/http"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (h *LoginHandlerImplementation) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get credentials
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Validate login form
	validationErrors, err := h.authUC.ValidateLoginForm(r.Context(), domain.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		http.Error(w, "Validation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(validationErrors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		html := "<div class='errors'>"
		for _, verr := range validationErrors {
			html += "<div class='error'>" + verr.Message + "</div>"
		}
		html += "</div>"
		w.Write([]byte(html))
		return
	}

	// Probe connection
	success, err := h.authUC.ProbeConnection(r.Context(), username, password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("<div class='error'>" + err.Error() + "</div>"))
		return
	}
	if !success {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("<div class='error'>Invalid credentials</div>"))
		return
	}

	// Get user accessible resources
	resources, err := h.authUC.GetUserAccessibleResources(r.Context(), username)
	if err != nil {
		http.Error(w, "Error getting accessible resources: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user has any accessible resources
	if len(resources.AccessibleDatabases) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<div class='error'>No accessible resources found</div>"))
		return
	}

	// Get first accessible database, schema, table
	database, err := h.authUC.GetFirstAccessibleDatabase(r.Context(), username)
	if err != nil {
		http.Error(w, "Error getting first database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	schema, err := h.authUC.GetFirstAccessibleSchema(r.Context(), username, database)
	if err != nil {
		http.Error(w, "Error getting first schema: "+err.Error(), http.StatusInternalServerError)
		return
	}

	table, err := h.authUC.GetFirstAccessibleTable(r.Context(), username, database, schema)
	if err != nil {
		http.Error(w, "Error getting first table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create session
	session, err := h.authUC.CreateSession(r.Context(), username, password, database, schema, table)
	if err != nil {
		http.Error(w, "Error creating session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
	})

	// Redirect to main view
	http.Redirect(w, r, "/main", http.StatusFound)
}
