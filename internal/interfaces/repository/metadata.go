package repository

import (
	"context"
	"github.com/kamil5b/lumen-pg/internal/domain"
)

// MetadataRepository handles PostgreSQL metadata operations
type MetadataRepository interface {
	// LoadGlobalMetadata fetches all databases, schemas, tables, columns, and relationships
	LoadGlobalMetadata(ctx context.Context) (*domain.GlobalMetadata, error)
	
	// LoadDatabaseMetadata fetches metadata for a specific database
	LoadDatabaseMetadata(ctx context.Context, dbName string) (*domain.DatabaseMetadata, error)
	
	// LoadTableMetadata fetches metadata for a specific table
	LoadTableMetadata(ctx context.Context, schemaName, tableName string) (*domain.TableMetadata, error)
	
	// LoadRoles fetches all PostgreSQL roles
	LoadRoles(ctx context.Context) ([]string, error)
	
	// LoadRolePermissions fetches accessible resources for a specific role
	LoadRolePermissions(ctx context.Context, roleName string) (*domain.RoleMetadata, error)
}
