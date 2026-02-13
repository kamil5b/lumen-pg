package repository

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// MetadataRepository defines operations for accessing cached database metadata
type MetadataRepository interface {
	// StoreMetadata stores metadata for a database
	StoreMetadata(ctx context.Context, metadata *domain.DatabaseMetadata) error

	// GetMetadata retrieves stored metadata for a database
	GetMetadata(ctx context.Context, database string) (*domain.DatabaseMetadata, error)

	// StoreRoleMetadata stores metadata for a PostgreSQL role
	StoreRoleMetadata(ctx context.Context, role string, metadata *domain.RoleMetadata) error

	// GetRoleMetadata retrieves stored metadata for a role
	GetRoleMetadata(ctx context.Context, role string) (*domain.RoleMetadata, error)

	// StoreAllRolesMetadata stores metadata for all roles at once
	StoreAllRolesMetadata(ctx context.Context, roles map[string]*domain.RoleMetadata) error

	// GetAllRolesMetadata retrieves metadata for all roles
	GetAllRolesMetadata(ctx context.Context) (map[string]*domain.RoleMetadata, error)

	// InvalidateMetadata clears cached metadata for a database
	InvalidateMetadata(ctx context.Context, database string) error

	// InvalidateRoleMetadata clears cached metadata for a role
	InvalidateRoleMetadata(ctx context.Context, role string) error

	// InvalidateAllMetadata clears all cached metadata
	InvalidateAllMetadata(ctx context.Context) error

	// GetAccessibleDatabases returns databases accessible by a role
	GetAccessibleDatabases(ctx context.Context, role string) ([]string, error)

	// GetAccessibleSchemas returns schemas accessible by a role in a database
	GetAccessibleSchemas(ctx context.Context, role, database string) ([]string, error)

	// GetAccessibleTables returns tables accessible by a role in a schema
	GetAccessibleTables(ctx context.Context, role, database, schema string) ([]string, error)

	// IsTableAccessible checks if a role can access a table
	IsTableAccessible(ctx context.Context, role, database, schema, table string) (bool, error)

	// GetTablePermissions returns the permissions a role has on a table
	GetTablePermissions(ctx context.Context, role, database, schema, table string) (*domain.AccessibleTable, error)
}
