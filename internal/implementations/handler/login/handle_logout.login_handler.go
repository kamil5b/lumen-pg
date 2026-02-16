package login

import (
	"net/http"
)

func (h *LoginHandlerImplementation) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Get session from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// No session, redirect to login
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Logout (invalidate session)
	err = h.authUC.Logout(r.Context(), cookie.Value)
	if err != nil {
		// Log error but still proceed with logout
		http.Error(w, "Error during logout: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusFound)
}
