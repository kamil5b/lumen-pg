package domain

// RoleMetadata represents a PostgreSQL role and its accessible resources
type RoleMetadata struct {
	RoleName            string
	AccessibleDatabases []string
	AccessibleSchemas   map[string][]string // database -> schemas
	AccessibleTables    map[string][]string // database.schema -> tables
	Permissions         map[string][]Permission
}

// Permission represents a database permission
type Permission string

const (
	PermissionConnect Permission = "CONNECT"
	PermissionSelect  Permission = "SELECT"
	PermissionInsert  Permission = "INSERT"
	PermissionUpdate  Permission = "UPDATE"
	PermissionDelete  Permission = "DELETE"
)
