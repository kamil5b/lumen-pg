package metadata

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// LoadGlobalMetadata fetches all databases, schemas, tables, columns, and relationships
func (m *MetadataRepository) LoadGlobalMetadata(ctx context.Context) (*domain.GlobalMetadata, error) {
	return nil, errors.New("not implemented yet")
}
