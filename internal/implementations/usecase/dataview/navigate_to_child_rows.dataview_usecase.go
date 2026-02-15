package dataview

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) NavigateToChildRows(ctx context.Context, username, database, schema, childTable, parentTable string, fkColumn, pkColumn, pkValue string) (*domain.QueryResult, error) {
	return nil, errors.New("not implemented")
}
