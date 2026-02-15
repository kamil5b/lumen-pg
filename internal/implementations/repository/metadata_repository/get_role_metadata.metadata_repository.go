package metadata_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (m *MetadataRepositoryImplementation) GetRoleMetadata(ctx context.Context, role string) (*domain.RoleMetadata, error) {
	return nil, errors.New("not implemented")
}
