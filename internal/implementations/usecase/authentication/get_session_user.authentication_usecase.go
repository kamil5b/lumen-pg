package authentication

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *AuthenticationUseCaseImplementation) GetSessionUser(ctx context.Context, sessionID string) (*domain.User, error) {
	return nil, errors.New("not implemented")
}
