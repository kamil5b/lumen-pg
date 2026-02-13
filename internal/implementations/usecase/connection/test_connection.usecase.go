package connection

import (
	"context"
	"errors"
)

// TestConnection tests a database connection string
func (c *ConnectionUseCase) TestConnection(ctx context.Context, connStr string) error {
	return errors.New("not implemented yet")
}
