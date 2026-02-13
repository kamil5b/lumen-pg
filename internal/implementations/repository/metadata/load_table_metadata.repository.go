package metadata

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// LoadTableMetadata fetches metadata for a specific table
func (m *MetadataRepository) LoadTableMetadata(ctx context.Context, schemaName, tableName string) (*domain.TableMetadata, error) {
	return nil, errors.New("not implemented yet")
}
