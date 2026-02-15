package query

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *QueryUseCaseImplementation) GetQueryAffectedRowCount(ctx context.Context, result *domain.QueryResult) int64 {
	if result == nil {
		return 0
	}
	return result.RowCount
}
