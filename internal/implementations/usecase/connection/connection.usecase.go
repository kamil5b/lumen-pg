package connection

// ConnectionUseCase is a noop implementation of the connection usecase
type ConnectionUseCase struct{}

// NewConnectionUseCase creates a new connection usecase
func NewConnectionUseCase() *ConnectionUseCase {
	return &ConnectionUseCase{}
}
