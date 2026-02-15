package metadata_repository

import (
	"context"
	"errors"
)

func (m *MetadataRepositoryImplementation) GetAccessibleDatabases(ctx context.Context, role string) ([]string, error) {
	return nil, errors.New("not implemented")
}
