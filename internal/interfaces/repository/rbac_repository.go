package repository

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// RBACRepository defines operations for role-based access control
type RBACRepository interface {
	// GetUserRole returns the role/username of an authenticated user
	GetUserRole(ctx context.Context, username string) (string, error)

	// GetAllRoles retrieves all PostgreSQL roles in the instance
	GetAllRoles(ctx context.Context) ([]string, error)

	// GetRolePermissions returns all permissions for a role on a specific table
	GetRolePermissions(ctx context.Context, role, database, schema, table string) (*domain.PermissionSet, error)

	// HasSelectPermission checks if a role can SELECT from a table
	HasSelectPermission(ctx context.Context, role, database, schema, table string) (bool, error)

	// HasInsertPermission checks if a role can INSERT into a table
	HasInsertPermission(ctx context.Context, role, database, schema, table string) (bool, error)

	// HasUpdatePermission checks if a role can UPDATE a table
	HasUpdatePermission(ctx context.Context, role, database, schema, table string) (bool, error)

	// HasDeletePermission checks if a role can DELETE from a table
	HasDeletePermission(ctx context.Context, role, database, schema, table string) (bool, error)

	// HasDatabaseConnectPermission checks if a role can CONNECT to a database
	HasDatabaseConnectPermission(ctx context.Context, role, database string) (bool, error)

	// HasSchemaUsagePermission checks if a role can USE a schema
	HasSchemaUsagePermission(ctx context.Context, role, database, schema string) (bool, error)

	// GetAccessibleDatabases returns all databases accessible by a role
	GetAccessibleDatabases(ctx context.Context, role string) ([]string, error)

	// GetAccessibleSchemas returns all schemas accessible by a role in a database
	GetAccessibleSchemas(ctx context.Context, role, database string) ([]string, error)

	// GetAccessibleTables returns all tables accessible by a role in a schema
	GetAccessibleTables(ctx context.Context, role, database, schema string) ([]domain.AccessibleTable, error)

	// CanAccessTable checks if a role can access a table
	CanAccessTable(ctx context.Context, role, database, schema, table string) (bool, error)

	// GetRoleMetadata returns complete metadata for a role including all accessible resources
	GetRoleMetadata(ctx context.Context, role string) (*domain.RoleMetadata, error)

	// IsReadOnlyRole checks if a role has read-only access (SELECT only)
	IsReadOnlyRole(ctx context.Context, role, database, schema, table string) (bool, error)

	// ValidateUserAccessToResource validates if a user has access to a specific resource
	ValidateUserAccessToResource(ctx context.Context, username, resourceType, database, schema, table string) (bool, error)
}
