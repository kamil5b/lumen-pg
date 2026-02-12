package domain

// DatabaseMetadata represents metadata for a single database
type DatabaseMetadata struct {
	Name    string
	Schemas []SchemaMetadata
}

// SchemaMetadata represents metadata for a single schema
type SchemaMetadata struct {
	Name   string
	Tables []TableMetadata
}

// TableMetadata represents metadata for a single table
type TableMetadata struct {
	SchemaName  string
	TableName   string
	Columns     []ColumnMetadata
	PrimaryKeys []string
	ForeignKeys []ForeignKeyMetadata
}

// ColumnMetadata represents metadata for a single column
type ColumnMetadata struct {
	Name         string
	DataType     string
	IsNullable   bool
	DefaultValue *string
}

// ForeignKeyMetadata represents a foreign key relationship
type ForeignKeyMetadata struct {
	ColumnName           string
	ReferencedSchemaName string
	ReferencedTableName  string
	ReferencedColumnName string
}

// GlobalMetadata represents all cached metadata for the PostgreSQL instance
type GlobalMetadata struct {
	Databases []DatabaseMetadata
	Roles     []RoleMetadata
}
