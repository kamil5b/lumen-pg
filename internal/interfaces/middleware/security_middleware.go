package middleware

import "net/http"

// SecurityMiddleware handles security concerns
type SecurityMiddleware interface {
	// SetSecurityHeaders sets appropriate security headers
	SetSecurityHeaders(next http.Handler) http.Handler

	// PreventCSRF prevents cross-site request forgery
	PreventCSRF(next http.Handler) http.Handler

	// RateLimiter limits request rate per user/IP
	RateLimiter(next http.Handler) http.Handler

	// ValidateCookieIntegrity validates that cookies haven't been tampered with
	ValidateCookieIntegrity(next http.Handler) http.Handler

	// EnforceSameSiteCookie enforces SameSite cookie attribute
	EnforceSameSiteCookie(next http.Handler) http.Handler

	// RequireHTTPSForCookies enforces HTTPS for secure cookies
	RequireHTTPSForCookies(next http.Handler) http.Handler
}
