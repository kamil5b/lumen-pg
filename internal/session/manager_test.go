package session_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/session"
	"github.com/stretchr/testify/assert"
)

// UC-S2-06: Session Cookie Creation - Username
func TestStubManager_CreateSessionCookies(t *testing.T) {
	mgr := session.NewStubManager([]byte("0123456789abcdef0123456789abcdef"))
	w := httptest.NewRecorder()
	err := mgr.CreateSessionCookies(w, "admin", "password123")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
}

// UC-S2-07: Session Cookie Creation - Password
func TestStubManager_CreateSessionCookies_EncryptsPassword(t *testing.T) {
	mgr := session.NewStubManager([]byte("0123456789abcdef0123456789abcdef"))
	w := httptest.NewRecorder()
	err := mgr.CreateSessionCookies(w, "admin", "password123")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
}

// UC-S2-08: Get Session From Cookies
func TestStubManager_GetSessionFromCookies(t *testing.T) {
	mgr := session.NewStubManager([]byte("0123456789abcdef0123456789abcdef"))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	s, err := mgr.GetSessionFromCookies(req)
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, s)
}

// UC-S2-12: Logout Cookie Clearing
func TestStubManager_ClearSessionCookies(t *testing.T) {
	mgr := session.NewStubManager([]byte("0123456789abcdef0123456789abcdef"))
	w := httptest.NewRecorder()
	// Should not panic even though it's a stub
	mgr.ClearSessionCookies(w)
}

// UC-S7-05: Cookie Tampering Detection
func TestStubManager_ValidateCookie(t *testing.T) {
	mgr := session.NewStubManager([]byte("0123456789abcdef0123456789abcdef"))
	cookie := &http.Cookie{
		Name:  "session",
		Value: "tampered-value",
	}
	err := mgr.ValidateCookie(cookie)
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
}

// UC-S7-06: Session Timeout Short-Lived Cookie
func TestStubManager_ShortLivedCookieTimeout(t *testing.T) {
	mgr := session.NewStubManager([]byte("0123456789abcdef0123456789abcdef"))
	// Verify the default short-lived timeout is set
	assert.NotZero(t, mgr.ShortLivedMaxAge)
}

// UC-S7-07: Session Timeout Long-Lived Cookie
func TestStubManager_LongLivedCookieTimeout(t *testing.T) {
	mgr := session.NewStubManager([]byte("0123456789abcdef0123456789abcdef"))
	// Verify the default long-lived timeout is set and longer than short-lived
	assert.NotZero(t, mgr.LongLivedMaxAge)
	assert.Greater(t, mgr.LongLivedMaxAge, mgr.ShortLivedMaxAge)
}
