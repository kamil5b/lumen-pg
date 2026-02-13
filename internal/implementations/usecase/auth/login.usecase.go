package auth

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// Login performs login operation
func (a *AuthUseCase) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	return nil, errors.New("not implemented yet")
}
