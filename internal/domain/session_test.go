package domain_test

import (
	"testing"
	"time"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/stretchr/testify/assert"
)

// UC-S2-01: Login Form Validation - Empty Username
func TestValidateLoginInput_EmptyUsername(t *testing.T) {
	err := domain.ValidateLoginInput("", "password123")
	assert.ErrorIs(t, err, domain.ErrEmptyUsername)
}

// UC-S2-02: Login Form Validation - Empty Password
func TestValidateLoginInput_EmptyPassword(t *testing.T) {
	err := domain.ValidateLoginInput("admin", "")
	assert.ErrorIs(t, err, domain.ErrEmptyPassword)
}

// UC-S2-01/02: Login Form Validation - Valid Input
func TestValidateLoginInput_Valid(t *testing.T) {
	err := domain.ValidateLoginInput("admin", "password123")
	assert.NoError(t, err)
}

// UC-S2-08: Session Validation - Valid Session
func TestSession_Validate_Valid(t *testing.T) {
	session := &domain.Session{
		Username:  "admin",
		Password:  "encrypted",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	err := session.Validate()
	assert.NoError(t, err)
}

// UC-S2-09: Session Validation - Expired Session
func TestSession_Validate_Expired(t *testing.T) {
	session := &domain.Session{
		Username:  "admin",
		Password:  "encrypted",
		CreatedAt: time.Now().Add(-1 * time.Hour),
		ExpiresAt: time.Now().Add(-30 * time.Minute),
	}
	err := session.Validate()
	assert.ErrorIs(t, err, domain.ErrSessionExpired)
}

// UC-S2-08: Session Validation - Invalid Session (empty username)
func TestSession_Validate_InvalidNoUsername(t *testing.T) {
	session := &domain.Session{
		Username:  "",
		Password:  "encrypted",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	err := session.Validate()
	assert.ErrorIs(t, err, domain.ErrInvalidSession)
}

// Session IsExpired
func TestSession_IsExpired(t *testing.T) {
	session := &domain.Session{
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}
	assert.True(t, session.IsExpired())
}

func TestSession_IsNotExpired(t *testing.T) {
	session := &domain.Session{
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	assert.False(t, session.IsExpired())
}
