package database_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (d *DatabaseRepositoryImplementation) GetTableData(ctx context.Context, params domain.TableDataParams) (*domain.QueryResult, error) {
	return nil, errors.New("not implemented")
}
