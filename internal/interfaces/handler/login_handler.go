package handler

import "net/http"

// LoginHandler handles authentication-related HTTP requests
type LoginHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	HandleLoginPage(w http.ResponseWriter, r *http.Request)
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleLogout(w http.ResponseWriter, r *http.Request)
}
