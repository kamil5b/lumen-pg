package erd

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *ERDUseCaseImplementation) GetERDTables(ctx context.Context, username, database, schema string) ([]domain.TableMetadata, error) {
	return nil, errors.New("not implemented")
}
