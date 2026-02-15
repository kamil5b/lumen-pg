package setup

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *SetupUseCaseImplementation) GetRoleAccessibility(ctx context.Context, role string) (*domain.RoleMetadata, error) {
	return nil, errors.New("not implemented")
}
