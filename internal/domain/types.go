package domain

import "time"

// DatabaseMetadata represents metadata about a database
type DatabaseMetadata struct {
	Name    string
	Schemas []SchemaMetadata
}

// SchemaMetadata represents metadata about a schema
type SchemaMetadata struct {
	Name   string
	Tables []TableMetadata
}

// TableMetadata represents metadata about a table
type TableMetadata struct {
	Name        string
	Columns     []ColumnMetadata
	PrimaryKeys []string
	ForeignKeys []ForeignKeyMetadata
}

// ColumnMetadata represents metadata about a column
type ColumnMetadata struct {
	Name       string
	DataType   string
	IsNullable bool
	IsPrimary  bool
}

// ForeignKeyMetadata represents metadata about a foreign key relationship
type ForeignKeyMetadata struct {
	ColumnName         string
	ReferencedTable    string
	ReferencedColumn   string
	ReferencedSchema   string
	ReferencedDatabase string
}

// RoleMetadata represents metadata about a PostgreSQL role
type RoleMetadata struct {
	Name                string
	AccessibleDatabases []string
	AccessibleSchemas   []string
	AccessibleTables    []AccessibleTable
}

// AccessibleTable represents a table accessible by a role
type AccessibleTable struct {
	Database  string
	Schema    string
	Name      string
	HasSelect bool
	HasInsert bool
	HasUpdate bool
	HasDelete bool
}

// User represents an authenticated user
type User struct {
	Username     string
	DatabaseName string
	SchemaName   string
	TableName    string
	ConnString   string
}

// Session represents a user session
type Session struct {
	ID        string
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// QueryResult represents the result of a SQL query execution
type QueryResult struct {
	Columns    []string
	Rows       []map[string]interface{}
	RowCount   int64
	TotalCount int64
	Error      string
}

// TransactionState represents an active transaction
type TransactionState struct {
	ID        string
	Username  string
	StartedAt time.Time
	ExpiresAt time.Time
	Edits     map[int]RowEdit
	Deletes   []int
	Inserts   []RowInsert
}

// RowEdit represents a buffered cell edit in a transaction
type RowEdit struct {
	RowIndex   int
	ColumnName string
	OldValue   interface{}
	NewValue   interface{}
}

// RowInsert represents a new row to be inserted
type RowInsert struct {
	Values map[string]interface{}
}

// QueryParams represents parameters for executing a query
type QueryParams struct {
	Query       string
	Offset      int
	Limit       int
	WhereClause string
	OrderBy     string
	OrderDir    string // ASC or DESC
}

// TableDataParams represents parameters for loading table data
type TableDataParams struct {
	Database    string
	Schema      string
	Table       string
	WhereClause string
	OrderBy     string
	OrderDir    string
	Offset      int
	Limit       int
	Cursor      string
}

// ForeignKeyInfo represents information about a foreign key relationship
type ForeignKeyInfo struct {
	ColumnName         string
	ReferencedTable    string
	ReferencedColumn   string
	ReferencedSchema   string
	ReferencedDatabase string
}

// ChildTableReference represents a reference from a parent table to a child table
type ChildTableReference struct {
	Database string
	Schema   string
	Table    string
	RowCount int64
}

// LoginRequest represents a login attempt
type LoginRequest struct {
	Username string
	Password string
}

// LoginResponse represents the result of a login attempt
type LoginResponse struct {
	Success   bool
	Message   string
	Username  string
	SessionID string
}

// CookieData represents encrypted data stored in a cookie
type CookieData struct {
	Username string
	Password string
	Nonce    string
}

// ConnectionString represents parsed connection string components
type ConnectionString struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	SSLMode  string
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface for ValidationError
func (ve ValidationError) Error() string {
	if ve.Field != "" {
		return ve.Field + ": " + ve.Message
	}
	return ve.Message
}
