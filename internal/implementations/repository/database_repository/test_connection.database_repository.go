package database_repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

func (d *DatabaseRepositoryImplementation) TestConnection(ctx context.Context, connString string) error {
	if connString == "" {
		return errors.New("connection string cannot be empty")
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}
