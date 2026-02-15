package authentication

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *AuthenticationUseCaseImplementation) CreateSession(ctx context.Context, username, password, database, schema, table string) (*domain.Session, error) {
	return nil, errors.New("not implemented")
}
