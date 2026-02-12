package domain

import (
	"errors"
	"time"
)

var (
	ErrSessionExpired     = errors.New("session expired")
	ErrInvalidSession     = errors.New("invalid session")
	ErrEmptyUsername       = errors.New("username cannot be empty")
	ErrEmptyPassword      = errors.New("password cannot be empty")
	ErrNoAccessibleResources = errors.New("no accessible resources found")
)

// Session represents a user session.
type Session struct {
	Username    string
	Password    string // encrypted
	CreatedAt   time.Time
	ExpiresAt   time.Time
	Database    string
	Schema      string
	Table       string
}

// IsExpired checks if the session has expired.
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// Validate checks if the session is valid.
func (s *Session) Validate() error {
	if s.Username == "" {
		return ErrInvalidSession
	}
	if s.IsExpired() {
		return ErrSessionExpired
	}
	return nil
}

// ValidateLoginInput validates login form input.
func ValidateLoginInput(username, password string) error {
	if username == "" {
		return ErrEmptyUsername
	}
	if password == "" {
		return ErrEmptyPassword
	}
	return nil
}
