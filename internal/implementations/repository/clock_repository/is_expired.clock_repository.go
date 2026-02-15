package clock_repository

import (
	"context"
	"time"
)

func (c *ClockRepositoryImplementation) IsExpired(ctx context.Context, expirationTime int64) bool {
	return time.Now().Unix() > expirationTime
}
