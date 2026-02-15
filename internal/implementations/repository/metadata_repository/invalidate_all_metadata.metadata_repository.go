package metadata_repository

import (
	"context"
	"errors"
)

func (m *MetadataRepositoryImplementation) InvalidateAllMetadata(ctx context.Context) error {
	return errors.New("not implemented")
}
