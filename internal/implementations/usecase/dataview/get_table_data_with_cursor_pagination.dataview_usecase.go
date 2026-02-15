package dataview

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) GetTableDataWithCursorPagination(ctx context.Context, username, database, schema, table, cursor string, limit int) (*domain.QueryResult, error) {
	return nil, errors.New("not implemented")
}
