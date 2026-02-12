package repository

import (
	"context"
	"github.com/kamil5b/lumen-pg/internal/domain"
)

// ConnectionRepository handles PostgreSQL connection operations
type ConnectionRepository interface {
	// ParseConnectionString parses a PostgreSQL connection string
	ParseConnectionString(connStr string) (*domain.Connection, error)
	
	// ValidateConnection tests connectivity to PostgreSQL
	ValidateConnection(ctx context.Context, conn *domain.Connection) error
	
	// Connect establishes a connection to PostgreSQL
	Connect(ctx context.Context, conn *domain.Connection) (interface{}, error)
}
