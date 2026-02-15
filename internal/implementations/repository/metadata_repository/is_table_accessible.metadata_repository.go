package metadata_repository

import (
	"context"
	"errors"
)

func (m *MetadataRepositoryImplementation) IsTableAccessible(ctx context.Context, role, database, schema, table string) (bool, error) {
	return false, errors.New("not implemented")
}
