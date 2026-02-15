package setup

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *SetupUseCaseImplementation) ParseConnectionString(ctx context.Context, connString string) (*domain.ConnectionString, error) {
	// Parse the connection string as a URL
	parsedURL, err := url.Parse(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Extract components from the parsed URL
	result := &domain.ConnectionString{
		Host:     parsedURL.Hostname(),
		Database: strings.TrimPrefix(parsedURL.Path, "/"),
	}

	// Extract port
	port := parsedURL.Port()
	if port != "" {
		result.Port = port
	} else {
		result.Port = "5432" // Default PostgreSQL port
	}

	// Extract username
	if parsedURL.User != nil {
		result.Username = parsedURL.User.Username()

		// Extract password
		if password, ok := parsedURL.User.Password(); ok {
			result.Password = password
		}
	}

	// Extract SSL mode from query parameters
	query := parsedURL.Query()
	if sslMode := query.Get("sslmode"); sslMode != "" {
		result.SSLMode = sslMode
	} else {
		result.SSLMode = "disable" // Default SSL mode
	}

	return result, nil
}
