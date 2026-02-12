package domain

// QueryRequest represents a SQL query execution request
type QueryRequest struct {
	SQL    string
	Params []interface{}
}

// QueryResult represents the result of a SQL query
type QueryResult struct {
	Columns      []string
	Rows         [][]interface{}
	RowsAffected int64
	TotalCount   int64 // Total count before pagination
	IsSelect     bool
	Error        error
}

// TableDataRequest represents a request for table data
type TableDataRequest struct {
	Schema      string
	Table       string
	WhereClause string
	OrderBy     string
	OrderDir    string // ASC or DESC
	Cursor      string
	Limit       int
}

// TableDataResult represents table data with pagination
type TableDataResult struct {
	Columns    []string
	Rows       [][]interface{}
	NextCursor string
	TotalCount int64
	HasMore    bool
}
