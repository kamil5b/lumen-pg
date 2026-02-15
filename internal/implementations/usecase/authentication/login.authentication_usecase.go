package authentication

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *AuthenticationUseCaseImplementation) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	return nil, errors.New("not implemented")
}
