package dataview

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) NavigateToParentRow(ctx context.Context, username, database, schema, table, columnName string, value interface{}) (*domain.QueryResult, error) {
	return nil, errors.New("not implemented")
}
