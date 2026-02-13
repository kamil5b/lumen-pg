package session

// SessionRepository is a noop implementation of the session repository
type SessionRepository struct{}

// NewSessionRepository creates a new session repository
func NewSessionRepository() *SessionRepository {
	return &SessionRepository{}
}
