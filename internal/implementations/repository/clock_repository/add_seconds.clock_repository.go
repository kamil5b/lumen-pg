package clock_repository

import (
	"context"
	"time"
)

func (c *ClockRepositoryImplementation) AddSeconds(ctx context.Context, seconds int64) int64 {
	return time.Now().Add(time.Duration(seconds) * time.Second).Unix()
}
