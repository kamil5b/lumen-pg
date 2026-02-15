package query

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *QueryUseCaseImplementation) ExecuteQueryWithPagination(ctx context.Context, username string, params domain.QueryParams) (*domain.QueryResult, error) {
	return nil, errors.New("not implemented")
}
