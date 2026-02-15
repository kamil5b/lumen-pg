package metadata_repository

import (
	"context"
	"errors"
)

func (m *MetadataRepositoryImplementation) GetAccessibleSchemas(ctx context.Context, role, database string) ([]string, error) {
	return nil, errors.New("not implemented")
}
