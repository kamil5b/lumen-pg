package clock_repository

import (
	"context"
	"time"
)

func (c *ClockRepositoryImplementation) NowNano(ctx context.Context) int64 {
	return time.Now().UnixNano()
}
