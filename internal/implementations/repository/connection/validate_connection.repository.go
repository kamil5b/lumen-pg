package connection

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// ValidateConnection tests connectivity to PostgreSQL
func (c *ConnectionRepository) ValidateConnection(ctx context.Context, conn *domain.Connection) error {
	return errors.New("not implemented yet")
}
