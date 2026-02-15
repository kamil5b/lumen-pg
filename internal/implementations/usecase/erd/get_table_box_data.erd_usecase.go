package erd

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *ERDUseCaseImplementation) GetTableBoxData(ctx context.Context, username, database, schema, table string) (*domain.TableMetadata, error) {
	return nil, errors.New("not implemented")
}
