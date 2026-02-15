package erd

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *ERDUseCaseImplementation) GetERDRelationships(ctx context.Context, username, database, schema string) ([]domain.ForeignKeyMetadata, error) {
	return nil, errors.New("not implemented")
}
