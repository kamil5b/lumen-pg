package domain

// DatabaseMetadata represents metadata about a database
type DatabaseMetadata struct {
	Name   string
	Tables []TableMetadata
}

// TableMetadata represents metadata about a table
type TableMetadata struct {
	Schema      string
	Name        string
	Columns     []ColumnMetadata
	PrimaryKeys []string
	ForeignKeys []ForeignKeyMetadata
}

// ColumnMetadata represents metadata about a column
type ColumnMetadata struct {
	Name     string
	DataType string
	Nullable bool
}

// ForeignKeyMetadata represents a foreign key relationship
type ForeignKeyMetadata struct {
	ColumnName       string
	ReferencedTable  string
	ReferencedSchema string
	ReferencedColumn string
}

// RolePermissions represents what resources a role can access
type RolePermissions struct {
	RoleName            string
	AccessibleDatabases []string
	AccessibleSchemas   map[string][]string   // database -> schemas
	AccessibleTables    map[string][]TableRef // database -> tables
}

// TableRef represents a reference to a table
type TableRef struct {
	Schema string
	Name   string
}

// GlobalMetadata represents all metadata cached in memory
type GlobalMetadata struct {
	Databases       []DatabaseMetadata
	RolePermissions map[string]*RolePermissions // roleName -> permissions
}
