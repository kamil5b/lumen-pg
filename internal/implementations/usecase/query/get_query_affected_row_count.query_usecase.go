package query

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *QueryUseCaseImplementation) GetQueryAffectedRowCount(ctx context.Context, result *domain.QueryResult) int64 {
	return 0
}
