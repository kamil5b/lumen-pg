package domain

const (
	DefaultPageSize = 50
	HardLimitRows   = 1000
)

// Cursor represents a pagination cursor for infinite scrolling.
type Cursor struct {
	LastValue  string
	LastID     int64
	PageSize   int
	TotalLoaded int
}

// NewCursor creates a new cursor with default page size.
func NewCursor() *Cursor {
	return &Cursor{
		PageSize: DefaultPageSize,
	}
}

// HasReachedLimit checks if the hard limit has been reached.
func (c *Cursor) HasReachedLimit() bool {
	return c.TotalLoaded >= HardLimitRows
}

// CanLoadMore checks if more data can be loaded.
func (c *Cursor) CanLoadMore() bool {
	return !c.HasReachedLimit()
}

// Advance advances the cursor after loading a page.
func (c *Cursor) Advance(lastValue string, lastID int64, loadedCount int) {
	c.LastValue = lastValue
	c.LastID = lastID
	c.TotalLoaded += loadedCount
}
