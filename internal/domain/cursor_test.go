package domain_test

import (
	"testing"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/stretchr/testify/assert"
)

// Cursor tests for Story 5 - Cursor Pagination

func TestCursor_NewCursor(t *testing.T) {
	c := domain.NewCursor()
	assert.Equal(t, domain.DefaultPageSize, c.PageSize)
	assert.Equal(t, 0, c.TotalLoaded)
	assert.False(t, c.HasReachedLimit())
	assert.True(t, c.CanLoadMore())
}

// UC-S5-08: Cursor Pagination Hard Limit
func TestCursor_HasReachedLimit(t *testing.T) {
	c := domain.NewCursor()
	c.TotalLoaded = domain.HardLimitRows
	assert.True(t, c.HasReachedLimit())
	assert.False(t, c.CanLoadMore())
}

// UC-S5-02: Cursor Pagination Next Page
func TestCursor_Advance(t *testing.T) {
	c := domain.NewCursor()
	c.Advance("value1", 50, 50)
	assert.Equal(t, "value1", c.LastValue)
	assert.Equal(t, int64(50), c.LastID)
	assert.Equal(t, 50, c.TotalLoaded)
	assert.True(t, c.CanLoadMore())
}

func TestCursor_Advance_MultiplePages(t *testing.T) {
	c := domain.NewCursor()
	for i := 0; i < 20; i++ {
		c.Advance("val", int64((i+1)*50), 50)
	}
	assert.Equal(t, 1000, c.TotalLoaded)
	assert.True(t, c.HasReachedLimit())
	assert.False(t, c.CanLoadMore())
}

func TestCursor_CanLoadMore_BelowLimit(t *testing.T) {
	c := domain.NewCursor()
	c.TotalLoaded = 500
	assert.True(t, c.CanLoadMore())
}

func TestCursor_CanLoadMore_AtLimit(t *testing.T) {
	c := domain.NewCursor()
	c.TotalLoaded = 1000
	assert.False(t, c.CanLoadMore())
}

func TestCursor_CanLoadMore_AboveLimit(t *testing.T) {
	c := domain.NewCursor()
	c.TotalLoaded = 1500
	assert.False(t, c.CanLoadMore())
}
