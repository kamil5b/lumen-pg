package metadata_repository

import (
	"context"
	"errors"
)

func (m *MetadataRepositoryImplementation) GetAccessibleTables(ctx context.Context, role, database, schema string) ([]string, error) {
	return nil, errors.New("not implemented")
}
