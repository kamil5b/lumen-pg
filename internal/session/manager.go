package session

import (
	"net/http"
	"time"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// Manager handles session cookie operations.
type Manager interface {
	// CreateSessionCookies creates both username and encrypted password cookies.
	CreateSessionCookies(w http.ResponseWriter, username, password string) error

	// GetSessionFromCookies extracts session information from cookies.
	GetSessionFromCookies(r *http.Request) (*domain.Session, error)

	// ClearSessionCookies clears all session cookies.
	ClearSessionCookies(w http.ResponseWriter)

	// ValidateCookie validates a cookie's integrity.
	ValidateCookie(cookie *http.Cookie) error
}

// StubManager is a stub implementation of Manager.
type StubManager struct {
	EncryptionKey     []byte
	ShortLivedMaxAge  time.Duration
	LongLivedMaxAge   time.Duration
}

func NewStubManager(encryptionKey []byte) *StubManager {
	return &StubManager{
		EncryptionKey:    encryptionKey,
		ShortLivedMaxAge: 30 * time.Minute,
		LongLivedMaxAge:  30 * 24 * time.Hour,
	}
}

func (s *StubManager) CreateSessionCookies(w http.ResponseWriter, username, password string) error {
	return domain.ErrNotImplemented
}

func (s *StubManager) GetSessionFromCookies(r *http.Request) (*domain.Session, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubManager) ClearSessionCookies(w http.ResponseWriter) {
	// stub - nothing to do
}

func (s *StubManager) ValidateCookie(cookie *http.Cookie) error {
	return domain.ErrNotImplemented
}
