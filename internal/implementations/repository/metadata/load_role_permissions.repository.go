package metadata

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// LoadRolePermissions fetches accessible resources for a specific role
func (m *MetadataRepository) LoadRolePermissions(ctx context.Context, roleName string) (*domain.RoleMetadata, error) {
	return nil, errors.New("not implemented yet")
}
