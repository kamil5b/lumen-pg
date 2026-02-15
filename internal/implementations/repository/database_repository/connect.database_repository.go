package database_repository

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
)

func (d *DatabaseRepositoryImplementation) Connect(ctx context.Context, connString string) error {
	if connString == "" {
		return fmt.Errorf("connection string cannot be empty")
	}

	if err := d.TestConnection(ctx, connString); err != nil {
		return err
	}

	return nil
}
