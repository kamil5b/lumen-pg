package domain

// ERDData represents entity-relationship diagram data
type ERDData struct {
	Tables        []ERDTable
	Relationships []ERDRelationship
}

// ERDTable represents a table in the ERD
type ERDTable struct {
	Schema  string
	Name    string
	Columns []ERDColumn
}

// ERDColumn represents a column in an ERD table
type ERDColumn struct {
	Name      string
	DataType  string
	IsPrimary bool
	IsForeign bool
}

// ERDRelationship represents a relationship between tables
type ERDRelationship struct {
	FromTable  TableRef
	ToTable    TableRef
	FromColumn string
	ToColumn   string
}
