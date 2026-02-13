package domain

// QueryRequest represents a request to execute a query
type QueryRequest struct {
	Query  string
	Offset int
	Limit  int
}

// QueryResponse represents a response from query execution
type QueryResponse struct {
	Success bool
	Data    *QueryResult
	Error   string
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Offset int
	Limit  int
	Cursor string
}

// FilterParams represents filter parameters
type FilterParams struct {
	WhereClause string
	OrderBy     string
	OrderDir    string
}

// SortParams represents sort parameters
type SortParams struct {
	Column    string
	Direction string
}

// DataExplorerNode represents a node in the data explorer tree
type DataExplorerNode struct {
	Type     string // "database", "schema", "table"
	Name     string
	Children []DataExplorerNode
}

// ERDTable represents a table in an ERD
type ERDTable struct {
	Name        string
	Columns     []ERDColumn
	X           int
	Y           int
	Width       int
	Height      int
	PrimaryKeys []string
}

// ERDColumn represents a column in an ERD table
type ERDColumn struct {
	Name       string
	DataType   string
	IsNullable bool
	IsPrimary  bool
}

// ERDRelationship represents a relationship in an ERD
type ERDRelationship struct {
	FromTable    string
	FromColumn   string
	ToTable      string
	ToColumn     string
	RelationType string // "one-to-many", "many-to-one", "one-to-one"
}

// ERDData represents complete ERD data
type ERDData struct {
	Tables        []ERDTable
	Relationships []ERDRelationship
	ZoomLevel     float64
	PanX          int
	PanY          int
}

// HeaderData represents data for the application header
type HeaderData struct {
	Username          string
	Database          string
	Schema            string
	Table             string
	ShowMetadataError bool
}

// TransactionStatus represents the status of a transaction
type TransactionStatus struct {
	IsActive      bool
	ID            string
	StartedAt     int64
	ExpiresAt     int64
	RemainingTime int64
	EditCount     int
	DeleteCount   int
	InsertCount   int
}

// PermissionSet represents a set of permissions
type PermissionSet struct {
	CanSelect  bool
	CanInsert  bool
	CanUpdate  bool
	CanDelete  bool
	CanConnect bool
	CanUsage   bool
}

// TableInfo represents information about a table
type TableInfo struct {
	Database    string
	Schema      string
	Name        string
	Columns     []ColumnInfo
	PrimaryKeys []string
	ForeignKeys []ForeignKeyInfo
	Permissions PermissionSet
	RowCount    int64
}

// ColumnInfo represents information about a column
type ColumnInfo struct {
	Name       string
	DataType   string
	IsNullable bool
	IsPrimary  bool
}

// DatabaseSelector represents available databases for selection
type DatabaseSelector struct {
	CurrentDatabase    string
	AvailableDatabases []string
}

// SchemaSelector represents available schemas for selection
type SchemaSelector struct {
	CurrentSchema    string
	AvailableSchemas []string
}

// TableSelector represents available tables for selection
type TableSelector struct {
	CurrentTable    string
	AvailableTables []string
}

// RowData represents a single row of data
type RowData struct {
	Values map[string]interface{}
	Index  int
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Type    string
	Message string
	Code    int
}

// SuccessResponse represents a successful response
type SuccessResponse struct {
	Message string
	Data    interface{}
}

// PaginationInfo represents pagination information
type PaginationInfo struct {
	CurrentOffset int
	CurrentLimit  int
	TotalRows     int64
	HasNext       bool
	HasPrevious   bool
	DisplayRows   int
}

// QueryExecutionMetrics represents metrics from query execution
type QueryExecutionMetrics struct {
	ExecutionTimeMs int64
	RowsAffected    int64
	RowsReturned    int64
}

// ConnectionTestResult represents the result of a connection test
type ConnectionTestResult struct {
	Success       bool
	Message       string
	ConnString    string
	FirstDatabase string
	FirstSchema   string
	FirstTable    string
}

// MetadataRefreshResult represents the result of metadata refresh
type MetadataRefreshResult struct {
	Success       bool
	Message       string
	DatabaseCount int
	SchemaCount   int
	TableCount    int
	RefreshTimeMs int64
}

// BulkOperationResult represents the result of a bulk operation
type BulkOperationResult struct {
	Success         bool
	AffectedRows    int64
	Errors          []string
	WarningCount    int
	ExecutionTimeMs int64
}

// DataExplorerItem represents a single item in the data explorer
type DataExplorerItem struct {
	Type     string // "database", "schema", "table"
	Name     string
	Children []DataExplorerItem
	IsRoot   bool
}

// TransactionEditBuffer represents all edits in a transaction
type TransactionEditBuffer struct {
	CellEdits     map[int]RowEdit
	Deletions     []int
	Insertions    []RowInsert
	FirstEditTime int64
	LastEditTime  int64
}

// ChildTableInfo represents information about a child table reference
type ChildTableInfo struct {
	Database         string
	Schema           string
	Table            string
	ForeignKeyColumn string
	RowCount         int64
}
