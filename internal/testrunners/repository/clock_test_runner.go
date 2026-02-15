package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// ClockRepositoryConstructor is a function type that creates a ClockRepository
type ClockRepositoryConstructor func() repository.ClockRepository

// ClockRepositoryRunner runs all clock repository tests against an implementation
// Maps to TEST_PLAN.md:
// - Story 2: Authentication & Identity [UC-S2-08, UC-S2-09: Session Validation with time checks]
// - Story 5: Main View & Data Interaction [UC-S5-14: Transaction Timer Expiration]
// - Story 7: Security & Best Practices [UC-S7-06~07: Session Timeout]
func ClockRepositoryRunner(t *testing.T, constructor ClockRepositoryConstructor) {
	t.Helper()

	ctx := context.Background()
	repo := constructor()

	// UC-S2-08: Session Validation - Valid Session
	// UC-S7-06: Session Timeout Short-Lived Cookie
	t.Run("Now returns current Unix timestamp", func(t *testing.T) {
		before := time.Now().Unix()
		now := repo.Now(ctx)
		after := time.Now().Unix()

		require.GreaterOrEqual(t, now, before)
		require.LessOrEqual(t, now, after+1)
	})

	// UC-S2-08: Session Validation - Valid Session
	t.Run("Now returns increasing timestamps", func(t *testing.T) {
		now1 := repo.Now(ctx)
		time.Sleep(10 * time.Millisecond)
		now2 := repo.Now(ctx)

		require.Greater(t, now2, now1)
	})

	// UC-S2-09: Session Validation - Expired Session
	// UC-S7-06: Session Timeout Short-Lived Cookie
	// IT-S7-03: Real Session Expiration
	t.Run("IsExpired returns true for past timestamp", func(t *testing.T) {
		pastTime := time.Now().Add(-1 * time.Hour).Unix()
		isExpired := repo.IsExpired(ctx, pastTime)

		require.True(t, isExpired)
	})

	// UC-S2-08: Session Validation - Valid Session
	// UC-S7-07: Session Timeout Long-Lived Cookie
	t.Run("IsExpired returns false for future timestamp", func(t *testing.T) {
		futureTime := time.Now().Add(1 * time.Hour).Unix()
		isExpired := repo.IsExpired(ctx, futureTime)

		require.False(t, isExpired)
	})

	// UC-S2-09: Session Validation - Expired Session
	t.Run("IsExpired returns true for very old timestamp", func(t *testing.T) {
		veryOldTime := int64(0)
		isExpired := repo.IsExpired(ctx, veryOldTime)

		require.True(t, isExpired)
	})

	// UC-S5-14: Transaction Timer Expiration
	// E2E-S7-04: Session Timeout Enforcement
	t.Run("IsExpired boundary test with current time", func(t *testing.T) {
		currentTime := time.Now().Unix()
		isExpired := repo.IsExpired(ctx, currentTime)

		// Current time should be considered expired (not in future)
		require.True(t, isExpired)
	})

	// UC-S7-06: Session Timeout Short-Lived Cookie
	t.Run("AddSeconds adds correct duration to current time", func(t *testing.T) {
		secondsToAdd := int64(3600)
		result := repo.AddSeconds(ctx, secondsToAdd)
		now := repo.Now(ctx)

		// Result should be approximately secondsToAdd seconds ahead of now
		require.Greater(t, result, now)
		require.Less(t, result-now, 5) // Allow small timing variance
	})

	// UC-S2-08: Session Validation - Valid Session
	t.Run("AddSeconds with zero returns approximately current time", func(t *testing.T) {
		result := repo.AddSeconds(ctx, 0)
		now := repo.Now(ctx)

		require.GreaterOrEqual(t, result, now)
		require.Less(t, result-now, 5)
	})

	// UC-S7-07: Session Timeout Long-Lived Cookie
	t.Run("AddSeconds with large duration", func(t *testing.T) {
		secondsToAdd := int64(86400 * 365) // One year
		result := repo.AddSeconds(ctx, secondsToAdd)
		now := repo.Now(ctx)

		require.GreaterOrEqual(t, result-now, secondsToAdd-5)
	})

	// UC-S5-14: Transaction Timer Expiration
	t.Run("TimeUntilExpiration returns positive for future timestamp", func(t *testing.T) {
		futureTime := time.Now().Add(1 * time.Hour).Unix()
		timeRemaining := repo.TimeUntilExpiration(ctx, futureTime)

		require.Greater(t, timeRemaining, int64(3595))
		require.Less(t, timeRemaining, int64(3605))
	})

	// UC-S2-09: Session Validation - Expired Session
	// UC-S5-14: Transaction Timer Expiration
	t.Run("TimeUntilExpiration returns zero or negative for past timestamp", func(t *testing.T) {
		pastTime := time.Now().Add(-1 * time.Hour).Unix()
		timeRemaining := repo.TimeUntilExpiration(ctx, pastTime)

		require.LessOrEqual(t, timeRemaining, int64(0))
	})

	// UC-S5-14: Transaction Timer Expiration
	t.Run("TimeUntilExpiration with very distant future", func(t *testing.T) {
		futureTime := time.Now().Add(24 * 30 * time.Hour).Unix() // ~30 days
		timeRemaining := repo.TimeUntilExpiration(ctx, futureTime)

		expectedSeconds := int64(30 * 24 * 3600)
		require.Greater(t, timeRemaining, expectedSeconds-10)
		require.Less(t, timeRemaining, expectedSeconds+10)
	})

	// UC-S2-08: Session Validation - Valid Session
	t.Run("NowNano returns current nanosecond timestamp", func(t *testing.T) {
		beforeNano := time.Now().UnixNano()
		nowNano := repo.NowNano(ctx)
		afterNano := time.Now().UnixNano()

		require.GreaterOrEqual(t, nowNano, beforeNano)
		require.LessOrEqual(t, nowNano, afterNano+1000000)
	})

	// UC-S2-08: Session Validation - Valid Session
	t.Run("NowNano returns increasing nanosecond timestamps", func(t *testing.T) {
		nowNano1 := repo.NowNano(ctx)
		time.Sleep(1 * time.Millisecond)
		nowNano2 := repo.NowNano(ctx)

		require.Greater(t, nowNano2, nowNano1)
	})

	// UC-S2-08: Session Validation - Valid Session
	t.Run("Now and NowNano are consistent", func(t *testing.T) {
		now := repo.Now(ctx)
		nowNano := repo.NowNano(ctx)

		// Convert to comparable units (both to nanoseconds)
		nowInNano := now * int64(time.Second)
		timeDiff := nowNano - nowInNano

		// Should be within a few milliseconds
		require.Greater(t, timeDiff, int64(-100000000)) // -100ms
		require.Less(t, timeDiff, int64(100000000))     // +100ms
	})

	// UC-S5-14: Transaction Timer Expiration
	// E2E-S5-11: Transaction Timer Countdown
	t.Run("Multiple TimeUntilExpiration calls for same timestamp decrease", func(t *testing.T) {
		futureTime := time.Now().Add(10 * time.Second).Unix()

		time1 := repo.TimeUntilExpiration(ctx, futureTime)
		time.Sleep(1 * time.Second)
		time2 := repo.TimeUntilExpiration(ctx, futureTime)

		require.Greater(t, time1, time2)
		require.GreaterOrEqual(t, time1-time2, int64(0))
		require.LessOrEqual(t, time1-time2, int64(2))
	})

	// UC-S7-06: Session Timeout Short-Lived Cookie
	t.Run("AddSeconds and IsExpired integration", func(t *testing.T) {
		expirationTime := repo.AddSeconds(ctx, 3600)
		isExpired := repo.IsExpired(ctx, expirationTime)

		require.False(t, isExpired)
	})

	// UC-S2-09: Session Validation - Expired Session
	t.Run("Session lifetime simulation with clock operations", func(t *testing.T) {
		// Create session expiration 1 hour in future
		sessionExpiry := repo.AddSeconds(ctx, 3600)

		// Check it's not expired
		require.False(t, repo.IsExpired(ctx, sessionExpiry))

		// Simulate checking time remaining
		timeRemaining := repo.TimeUntilExpiration(ctx, sessionExpiry)
		require.Greater(t, timeRemaining, int64(3595))

		// Create an old timestamp and check it's expired
		oldTime := repo.AddSeconds(ctx, -3600)
		require.True(t, repo.IsExpired(ctx, oldTime))
	})
}
