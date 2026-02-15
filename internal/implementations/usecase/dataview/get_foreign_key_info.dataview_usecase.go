package dataview

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) GetForeignKeyInfo(ctx context.Context, username, database, schema, table string) ([]domain.ForeignKeyInfo, error) {
	return nil, errors.New("not implemented")
}
