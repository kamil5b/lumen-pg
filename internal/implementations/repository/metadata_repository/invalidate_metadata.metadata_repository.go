package metadata_repository

import (
	"context"
	"errors"
)

func (m *MetadataRepositoryImplementation) InvalidateMetadata(ctx context.Context, database string) error {
	return errors.New("not implemented")
}
