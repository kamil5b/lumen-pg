package auth

// AuthUseCase is a noop implementation of the auth usecase
type AuthUseCase struct{}

// NewAuthUseCase creates a new auth usecase
func NewAuthUseCase() *AuthUseCase {
	return &AuthUseCase{}
}
