package usecase

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// RBACUseCase defines operations for role-based access control
type RBACUseCase interface {
	// CheckTableAccess checks if a user has access to a table
	CheckTableAccess(ctx context.Context, username, database, schema, table string) (bool, error)

	// CheckSelectPermission checks if a user can SELECT from a table
	CheckSelectPermission(ctx context.Context, username, database, schema, table string) (bool, error)

	// CheckInsertPermission checks if a user can INSERT into a table
	CheckInsertPermission(ctx context.Context, username, database, schema, table string) (bool, error)

	// CheckUpdatePermission checks if a user can UPDATE a table
	CheckUpdatePermission(ctx context.Context, username, database, schema, table string) (bool, error)

	// CheckDeletePermission checks if a user can DELETE from a table
	CheckDeletePermission(ctx context.Context, username, database, schema, table string) (bool, error)

	// CheckDatabaseAccess checks if a user can access a database
	CheckDatabaseAccess(ctx context.Context, username, database string) (bool, error)

	// CheckSchemaAccess checks if a user can access a schema in a database
	CheckSchemaAccess(ctx context.Context, username, database, schema string) (bool, error)

	// GetUserAccessibleDatabases returns all databases accessible by a user
	GetUserAccessibleDatabases(ctx context.Context, username string) ([]string, error)

	// GetUserAccessibleSchemas returns all schemas accessible by a user in a database
	GetUserAccessibleSchemas(ctx context.Context, username, database string) ([]string, error)

	// GetUserAccessibleTables returns all tables accessible by a user in a schema
	GetUserAccessibleTables(ctx context.Context, username, database, schema string) ([]domain.AccessibleTable, error)

	// GetTablePermissions returns the permissions a user has on a table
	GetTablePermissions(ctx context.Context, username, database, schema, table string) (*domain.PermissionSet, error)

	// IsTableReadOnly checks if a table is read-only for the user
	IsTableReadOnly(ctx context.Context, username, database, schema, table string) (bool, error)

	// ValidateUserAccessToResource validates if a user has access to a specific resource
	ValidateUserAccessToResource(ctx context.Context, username, resourceType, database, schema, table string) (bool, error)

	// GetUserRole returns the role of a user
	GetUserRole(ctx context.Context, username string) (string, error)

	// VerifyUserPermissions verifies all permissions for a user on a table
	VerifyUserPermissions(ctx context.Context, username, database, schema, table string) (*domain.PermissionSet, error)
}
