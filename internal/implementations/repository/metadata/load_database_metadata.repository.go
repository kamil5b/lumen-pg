package metadata

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// LoadDatabaseMetadata fetches metadata for a specific database
func (m *MetadataRepository) LoadDatabaseMetadata(ctx context.Context, dbName string) (*domain.DatabaseMetadata, error) {
	return nil, errors.New("not implemented yet")
}
