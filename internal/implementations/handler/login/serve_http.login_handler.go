package login

import "net/http"

func (h *LoginHandlerImplementation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login":
		if r.Method == http.MethodGet {
			h.HandleLoginPage(w, r)
		} else if r.Method == http.MethodPost {
			h.HandleLogin(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/logout":
		h.HandleLogout(w, r)
	default:
		http.NotFound(w, r)
	}
}
