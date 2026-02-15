package setup

import (
	"context"
	"fmt"
)

func (u *SetupUseCaseImplementation) TestSuperadminConnection(ctx context.Context, connString string) (bool, error) {
	// Test the connection using the database repository
	err := u.databaseRepo.TestConnection(ctx, connString)
	if err != nil {
		return false, fmt.Errorf("failed to test superadmin connection: %w", err)
	}

	return true, nil
}
