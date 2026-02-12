package domain

import "time"

// Column represents a database column.
type Column struct {
	Name       string
	DataType   string
	IsNullable bool
	IsPrimary  bool
	Default    string
}

// ForeignKey represents a foreign key relationship.
type ForeignKey struct {
	ConstraintName    string
	SourceTable       string
	SourceSchema      string
	SourceColumn      string
	TargetTable       string
	TargetSchema      string
	TargetColumn      string
	SourceDatabase    string
	TargetDatabase    string
}

// Table represents a database table.
type Table struct {
	Name        string
	Schema      string
	Database    string
	Columns     []Column
	ForeignKeys []ForeignKey
}

// Schema represents a database schema.
type Schema struct {
	Name     string
	Database string
	Tables   []Table
}

// Database represents a PostgreSQL database.
type Database struct {
	Name    string
	Schemas []Schema
}

// Role represents a PostgreSQL role with permissions.
type Role struct {
	Name                string
	AccessibleDatabases []string
	AccessibleSchemas   map[string][]string // db -> schemas
	AccessibleTables    map[string][]string // "db.schema" -> tables
}

// Metadata holds all cached metadata from the PostgreSQL instance.
type Metadata struct {
	Databases []Database
	Roles     []Role
	UpdatedAt time.Time
}

// ERDData represents Entity-Relationship Diagram data.
type ERDData struct {
	Tables        []ERDTable
	Relationships []ERDRelationship
}

// ERDTable represents a table in the ERD.
type ERDTable struct {
	Name    string
	Columns []ERDColumn
}

// ERDColumn represents a column in the ERD.
type ERDColumn struct {
	Name     string
	DataType string
	IsPK     bool
	IsFK     bool
}

// ERDRelationship represents a relationship in the ERD.
type ERDRelationship struct {
	SourceTable  string
	SourceColumn string
	TargetTable  string
	TargetColumn string
}

// QueryResult represents the result of a query execution.
type QueryResult struct {
	Columns      []string
	Rows         [][]interface{}
	AffectedRows int64
	Message      string
	TotalSize    int64
	IsSelect     bool
}

// ReferencingTable represents a table that references a primary key.
type ReferencingTable struct {
	TableName string
	Schema    string
	Database  string
	RowCount  int64
}

// CursorPage represents a page of cursor-paginated data.
type CursorPage struct {
	Rows       [][]interface{}
	Columns    []string
	NextCursor string
	HasMore    bool
	TotalSize  int64
}

// ConnectionConfig holds parsed connection string parameters.
type ConnectionConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
	SSLMode  string
}
