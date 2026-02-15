package dataview

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) SortTableData(ctx context.Context, username, database, schema, table, orderBy, orderDir string, offset, limit int) (*domain.QueryResult, error) {
	return nil, errors.New("not implemented")
}
