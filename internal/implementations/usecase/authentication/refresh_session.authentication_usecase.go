package authentication

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *AuthenticationUseCaseImplementation) RefreshSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	return nil, errors.New("not implemented")
}
