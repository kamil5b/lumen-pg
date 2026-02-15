package metadata_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (m *MetadataRepositoryImplementation) StoreAllRolesMetadata(ctx context.Context, roles map[string]*domain.RoleMetadata) error {
	return errors.New("not implemented")
}
