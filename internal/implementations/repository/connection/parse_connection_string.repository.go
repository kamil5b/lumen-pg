package connection

import (
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// ParseConnectionString parses a PostgreSQL connection string
func (c *ConnectionRepository) ParseConnectionString(connStr string) (*domain.Connection, error) {
	return nil, errors.New("not implemented yet")
}
