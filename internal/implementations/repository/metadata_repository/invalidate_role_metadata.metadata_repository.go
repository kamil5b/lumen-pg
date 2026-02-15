package metadata_repository

import (
	"context"
	"errors"
)

func (m *MetadataRepositoryImplementation) InvalidateRoleMetadata(ctx context.Context, role string) error {
	return errors.New("not implemented")
}
