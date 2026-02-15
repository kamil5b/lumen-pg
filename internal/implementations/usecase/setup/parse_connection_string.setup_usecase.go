package setup

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *SetupUseCaseImplementation) ParseConnectionString(ctx context.Context, connString string) (*domain.ConnectionString, error) {
	return nil, errors.New("not implemented")
}
