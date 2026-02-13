package metadata

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// LoadRoleAccessibleResources loads accessible resources for a role
func (m *MetadataUseCase) LoadRoleAccessibleResources(ctx context.Context, roleName string) (*domain.RoleMetadata, error) {
	return nil, errors.New("not implemented yet")
}
