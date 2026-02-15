package authentication

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *AuthenticationUseCaseImplementation) GetUserAccessibleResources(ctx context.Context, username string) (*domain.RoleMetadata, error) {
	return nil, errors.New("not implemented")
}
