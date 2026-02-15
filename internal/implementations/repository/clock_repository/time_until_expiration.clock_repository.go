package clock_repository

import (
	"context"
)

func (c *ClockRepositoryImplementation) TimeUntilExpiration(ctx context.Context, expirationTime int64) int64 {
	remaining := expirationTime - c.Now(ctx)
	if remaining < 0 {
		return 0
	}
	return remaining
}
