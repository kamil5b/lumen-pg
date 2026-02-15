package authentication

import (
	"context"
	"fmt"
)

func (u *AuthenticationUseCaseImplementation) ProbeConnection(ctx context.Context, username, password string) (bool, error) {
	// Build a connection string from username and password
	// Use default host and port
	connString := fmt.Sprintf("postgres://%s:%s@localhost:5432/postgres?sslmode=disable", username, password)

	err := u.databaseRepo.TestConnection(ctx, connString)
	if err != nil {
		return false, err
	}
	return true, nil
}
