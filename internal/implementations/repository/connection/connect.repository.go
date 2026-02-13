package connection

import (
	"context"
	"errors"
)

// Connect establishes a connection to PostgreSQL
func (c *ConnectionRepository) Connect(ctx context.Context, conn interface{}) (interface{}, error) {
	return nil, errors.New("not implemented yet")
}
