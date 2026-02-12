package domain

import "time"

// Session represents a user's authenticated session
type Session struct {
	Username          string
	EncryptedPassword string
	FirstAccessibleDB string
	CreatedAt         time.Time
	ExpiresAt         time.Time
}

// SessionCookies represents the cookies for session management
type SessionCookies struct {
	UsernameToken string // Long-lived identity cookie
	SessionToken  string // Short-lived encrypted password cookie
}
