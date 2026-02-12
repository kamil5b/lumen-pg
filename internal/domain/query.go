package domain

// QueryRequest represents a query execution request
type QueryRequest struct {
	SQL        string
	Parameters []interface{}
}

// QueryResult represents the result of a query execution
type QueryResult struct {
	Columns      []string
	Rows         [][]interface{}
	TotalRows    int64 // Total rows available (for pagination)
	LoadedRows   int   // Rows actually loaded (max 1000)
	AffectedRows int64 // For DML queries
	Success      bool
	ErrorMessage string
}

// PaginationCursor represents pagination state for cursor-based pagination
type PaginationCursor struct {
	LastID      interface{}
	Offset      int
	Limit       int
	HasMore     bool
	TotalCount  int64
	LoadedCount int
}
