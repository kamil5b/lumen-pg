package usecase

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// SetupUseCase defines operations for application setup and configuration
type SetupUseCase interface {
	// ValidateConnectionString validates a PostgreSQL connection string format
	ValidateConnectionString(ctx context.Context, connString string) (bool, error)

	// ParseConnectionString parses a connection string into components
	ParseConnectionString(ctx context.Context, connString string) (*domain.ConnectionString, error)

	// TestSuperadminConnection tests connectivity using superadmin credentials
	TestSuperadminConnection(ctx context.Context, connString string) (bool, error)

	// InitializeMetadata fetches and caches metadata for all databases, schemas, tables, and relations
	InitializeMetadata(ctx context.Context, connString string) error

	// InitializeRBAC fetches and caches RBAC metadata for all PostgreSQL roles
	InitializeRBAC(ctx context.Context, connString string) error

	// GetAllRoles retrieves all PostgreSQL roles in the instance
	GetAllRoles(ctx context.Context) ([]string, error)

	// GetRoleAccessibility determines which resources are accessible by a role
	GetRoleAccessibility(ctx context.Context, role string) (*domain.RoleMetadata, error)

	// RefreshMetadata reloads all cached metadata from the database
	RefreshMetadata(ctx context.Context) error

	// RefreshRBACMetadata reloads all cached RBAC metadata from the database
	RefreshRBACMetadata(ctx context.Context) error

	// IsInitialized checks if the application has been properly initialized
	IsInitialized(ctx context.Context) (bool, error)
}
