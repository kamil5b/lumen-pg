package metadata

import (
	"context"
	"errors"
)

// LoadRoles fetches all PostgreSQL roles
func (m *MetadataRepository) LoadRoles(ctx context.Context) ([]string, error) {
	return nil, errors.New("not implemented yet")
}
