package connection

// ConnectionRepository is a noop implementation of the connection repository
type ConnectionRepository struct{}

// NewConnectionRepository creates a new connection repository
func NewConnectionRepository() *ConnectionRepository {
	return &ConnectionRepository{}
}
