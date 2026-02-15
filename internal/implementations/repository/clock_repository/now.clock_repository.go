package clock_repository

import (
	"context"
	"time"
)

func (c *ClockRepositoryImplementation) Now(ctx context.Context) int64 {
	return time.Now().Unix()
}
