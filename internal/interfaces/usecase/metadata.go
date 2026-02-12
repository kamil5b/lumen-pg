package usecase

import (
	"context"
	"github.com/kamil5b/lumen-pg/internal/domain"
)

// MetadataUseCase handles metadata loading operations
type MetadataUseCase interface {
	LoadGlobalMetadata(ctx context.Context) (*domain.GlobalMetadata, error)
	LoadRoleAccessibleResources(ctx context.Context, roleName string) (*domain.RoleMetadata, error)
}
