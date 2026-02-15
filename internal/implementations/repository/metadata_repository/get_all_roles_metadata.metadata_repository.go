package metadata_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (m *MetadataRepositoryImplementation) GetAllRolesMetadata(ctx context.Context) (map[string]*domain.RoleMetadata, error) {
	return nil, errors.New("not implemented")
}
