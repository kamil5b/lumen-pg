package metadata_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (m *MetadataRepositoryImplementation) GetTablePermissions(ctx context.Context, role, database, schema, table string) (*domain.AccessibleTable, error) {
	return nil, errors.New("not implemented")
}
