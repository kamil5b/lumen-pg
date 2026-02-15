package authentication

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *AuthenticationUseCaseImplementation) GetUserAccessibleResources(ctx context.Context, username string) (*domain.RoleMetadata, error) {
	return u.metadataRepo.GetRoleMetadata(ctx, username)
}
