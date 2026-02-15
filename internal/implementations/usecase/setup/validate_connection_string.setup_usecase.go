package setup

import (
	"context"
	"net/url"
	"strings"
)

func (u *SetupUseCaseImplementation) ValidateConnectionString(ctx context.Context, connString string) (bool, error) {
	// Validate that the connection string is not empty
	if strings.TrimSpace(connString) == "" {
		return false, nil
	}

	// Try to parse the connection string as a URL
	// PostgreSQL connection strings can be in the format: postgres://user:password@host:port/database?sslmode=disable
	parsedURL, err := url.Parse(connString)
	if err != nil {
		return false, nil
	}

	// Check if the scheme is postgres or postgresql
	if parsedURL.Scheme != "postgres" && parsedURL.Scheme != "postgresql" {
		return false, nil
	}

	// Check if host is present
	if parsedURL.Host == "" {
		return false, nil
	}

	// Connection string is valid
	return true, nil
}
