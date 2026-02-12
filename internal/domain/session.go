package domain

import "time"

// Session represents a user session
type Session struct {
	Username           string
	EncryptedPassword  string
	AccessibleMetadata *RoleMetadata
	CreatedAt          time.Time
	ExpiresAt          time.Time
}

// LoginRequest represents a login attempt
type LoginRequest struct {
	Username string
	Password string
}

// LoginResponse represents the result of a login attempt
type LoginResponse struct {
	Success            bool
	ErrorMessage       string
	Session            *Session
	FirstAccessibleDB  string
	FirstAccessibleTbl string
}
