package metadata

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// LoadGlobalMetadata loads global metadata
func (m *MetadataUseCase) LoadGlobalMetadata(ctx context.Context) (*domain.GlobalMetadata, error) {
	return nil, errors.New("not implemented yet")
}
