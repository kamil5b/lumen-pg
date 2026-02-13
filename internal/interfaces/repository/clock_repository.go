package repository

import (
	"context"
)

// ClockRepository defines operations for time-related functionality
type ClockRepository interface {
	// Now returns the current time in Unix seconds
	Now(ctx context.Context) int64

	// NowNano returns the current time in nanoseconds
	NowNano(ctx context.Context) int64

	// AddSeconds adds seconds to the current time and returns Unix timestamp
	AddSeconds(ctx context.Context, seconds int64) int64

	// IsExpired checks if a given Unix timestamp is in the past
	IsExpired(ctx context.Context, expirationTime int64) bool

	// TimeUntilExpiration returns seconds until a given expiration time
	TimeUntilExpiration(ctx context.Context, expirationTime int64) int64
}
