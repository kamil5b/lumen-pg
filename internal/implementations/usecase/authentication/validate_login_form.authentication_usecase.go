package authentication

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *AuthenticationUseCaseImplementation) ValidateLoginForm(ctx context.Context, req domain.LoginRequest) ([]domain.ValidationError, error) {
	return nil, errors.New("not implemented")
}
