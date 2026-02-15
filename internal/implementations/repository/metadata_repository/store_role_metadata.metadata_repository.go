package metadata_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (m *MetadataRepositoryImplementation) StoreRoleMetadata(ctx context.Context, role string, metadata *domain.RoleMetadata) error {
	return errors.New("not implemented")
}
