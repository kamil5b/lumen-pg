package metadata_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (m *MetadataRepositoryImplementation) GetMetadata(ctx context.Context, database string) (*domain.DatabaseMetadata, error) {
	return nil, errors.New("not implemented")
}
