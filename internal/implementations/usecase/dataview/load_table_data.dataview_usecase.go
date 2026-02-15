package dataview

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) LoadTableData(ctx context.Context, username string, params domain.TableDataParams) (*domain.QueryResult, error) {
	return nil, errors.New("not implemented")
}
